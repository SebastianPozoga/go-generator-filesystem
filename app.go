package main

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"unicode"
	"unicode/utf8"

	"github.com/SebastianPozoga/go-generator-filesystem/cache"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/filespace/diskfs"
	"github.com/goatcms/goatcore/filesystem/fsloop"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// App is main application object. Contains all app data
type App struct {
	From, To, Cache           string
	FS, FromFS, ToFS, CacheFS filesystem.Filespace
	modTimes                  *cache.ModTimes
}

func (app *App) Valid() {
	if app.From == "" {
		panic("from flag is required")
	}
	if app.To == "" {
		panic("to flag is required")
	}
}

func (app *App) InitFS() {
	var err error
	if app.FromFS, err = diskfs.NewFilespace(app.From); err != nil {
		panic(err)
	}
	if app.ToFS, err = diskfs.NewFilespace(app.To); err != nil {
		panic(err)
	}
	if err = app.ToFS.RemoveAll("."); err != nil {
		panic(err)
	}
	if err = app.ToFS.MkdirAll(".", filesystem.DefaultUnixDirMode); err != nil {
		panic(err)
	}
	if app.Cache != "" {
		if app.CacheFS, err = diskfs.NewFilespace(app.Cache); err != nil {
			panic(err)
		}
		if err = app.CacheFS.MkdirAll(".", filesystem.DefaultUnixDirMode); err != nil {
			panic(err)
		}
	}
}

func (app *App) readModTimes() (err error) {
	var reader filesystem.Reader
	app.modTimes = cache.NewModTimes()
	if app.CacheFS == nil || !app.CacheFS.IsFile(CacheModPath) {
		return
	}
	if reader, err = app.CacheFS.Reader(CacheModPath); err != nil {
		return
	}
	defer reader.Close()
	err = app.modTimes.Read(reader)
	return
}

func (app *App) persistModTimes() (err error) {
	var writer filesystem.Writer
	if app.CacheFS == nil {
		return
	}
	app.CacheFS.Remove(CacheModPath)
	if writer, err = app.CacheFS.Writer(CacheModPath); err != nil {
		return
	}
	defer writer.Close()
	err = app.modTimes.Write(writer)
	return
}

func (app *App) Run() (err error) {
	if err = app.readModTimes(); err != nil {
		return
	}
	loop := fsloop.NewLoop(&fsloop.LoopData{
		Filespace: app.FromFS,
		OnDir: func(fs filesystem.Filespace, subPath string) error {
			return fs.MkdirAll(subPath, filesystem.DefaultUnixDirMode)
		},
		FileFilter: func(fs filesystem.Filespace, subPath string) (modified bool) {
			modified = app.modTimes.IsFileModified(fs, subPath)
			if !modified {
				fmt.Printf("\n [no modified] %s", subPath)
			}
			return
		},
		OnFile: func(fs filesystem.Filespace, subPath string) (err error) {
			var (
				bytes       []byte
				checksum    []byte
				filename    = toCamelCase(filepath.Base(subPath), true)
				dirname     = toUnderscore(filepath.Base(filepath.Dir(subPath)))
				lstat       os.FileInfo
				contentType string
			)
			if bytes, err = fs.ReadFile(subPath); err != nil {
				panic(err)
			}
			if checksum, err = calculateChecksum(bytes); err != nil {
				panic(err)
			}
			if dirname == "" {
				dirname = toUnderscore(filepath.Base(app.To))
			}
			if dirname == "" {
				dirname = "fs"
			}
			contentType = mime.TypeByExtension(filepath.Ext(subPath))
			if contentType == "" {
				contentType = http.DetectContentType(bytes)
			}
			if lstat, err = fs.Lstat(subPath); err != nil {
				panic(err)
			}
			file := &GOBinaryFile{
				Bytes:       bytes,
				Checksum:    checksum,
				Package:     dirname,
				Name:        filename,
				ModTime:     lstat.ModTime(),
				ContentType: contentType,
			}
			firstLetter, _ := utf8.DecodeRuneInString(file.Package)
			if !unicode.IsLetter(firstLetter) {
				panic(fmt.Sprintf("Package name must start from letter - your directory name is not start from letter (%s)", subPath))
			}
			result := file.String()
			destPath := subPath + ".go"
			if err = app.ToFS.WriteFile(destPath, []byte(result), filesystem.DefaultUnixFileMode); err != nil {
				panic(err)
			}
			app.modTimes.Add(cache.ModTime{
				Path:    subPath,
				ModTime: lstat.ModTime(),
			})
			fmt.Printf("\n [generated] %s", destPath)
			return
		},
		Consumers:  runtime.NumCPU(),
		Producents: 1,
	}, nil)
	loop.Run("./")
	loop.Wait()
	fmt.Println()
	if err = app.persistModTimes(); err != nil {
		return
	}
	return goaterr.ToError(loop.Errors())
}
