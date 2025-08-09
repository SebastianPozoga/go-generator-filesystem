package main

import (
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"unicode"
	"unicode/utf8"

	"github.com/SebastianPozoga/go-generator-filesystem/cache"
	"github.com/SebastianPozoga/go-generator-filesystem/formatters/goformatter"
	"github.com/SebastianPozoga/go-generator-filesystem/names"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/filespace/diskfs"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

type processFileRow struct {
	modTimes *cache.ModTimes
	mapFile  *goformatter.MapFile
}

// App is main application object. Contains all app data
type App struct {
	PackagePrefix             string
	From, To, Cache           string
	FS, FromFS, ToFS, CacheFS filesystem.Filespace
	LogAll                    bool
	modTimes                  *cache.ModTimes

	// errorsMu sync.Mutex
	errors []error

	isModified, isRemoved, isCreated bool
}

func (app *App) Valid() {
	if app.From == "" {
		panic("from flag is required")
	}
	if app.To == "" {
		panic("to flag is required")
	}
}

func (app *App) ExLog(msg string, args ...any) {
	if app.LogAll {
		fmt.Printf(msg, args...)
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
	// if err = app.ToFS.RemoveAll("."); err != nil {
	// 	panic(err)
	// }
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

func (app *App) persistModTimes(modTimes *cache.ModTimes) (err error) {
	var writer filesystem.Writer
	if app.CacheFS == nil {
		return
	}
	app.CacheFS.Remove(CacheModPath)
	if writer, err = app.CacheFS.Writer(CacheModPath); err != nil {
		return
	}
	defer writer.Close()
	err = modTimes.Write(writer)
	return
}

func (app *App) persistMap(m *goformatter.MapFile) (err error) {
	if app.ToFS == nil {
		return
	}
	app.ToFS.Remove(MapFilePath)
	err = app.ToFS.WriteFile(MapFilePath, []byte(m.String()), filesystem.DefaultUnixFileMode)
	return
}

func (app *App) readDir(result processFileRow, dirPath string, changedFileChan chan names.FileNames) (err error) {
	var (
		infos          []fs.FileInfo
		basePath       = dirPath + "/"
		defaultDirName = names.ToUnderscore(filepath.Base(app.To))
	)
	if defaultDirName == "" {
		defaultDirName = "fs"
	}
	if infos, err = app.FromFS.ReadDir(dirPath); err != nil {
		return
	}
	for _, node := range infos {
		if node.Name() == "." || node.Name() == ".." || node.Name() == ".DsStore" {
			continue
		}
		nodePath := basePath + node.Name()
		if node.IsDir() {
			app.ToFS.MkdirAll(nodePath, filesystem.DefaultUnixDirMode)
			app.readDir(result, nodePath, changedFileChan)
		} else {
			fileNames := names.NewFileNames(nodePath, defaultDirName)
			modified, creted := app.modTimes.IsFileModified(app.FromFS, fileNames.Path)
			app.isModified = app.isModified || modified
			app.isCreated = app.isCreated || creted

			if modified {
				app.ExLog("\n [modified] %s", fileNames.Path)
			} else if creted {
				app.ExLog("\n [created] %s", fileNames.Path)
			} else {
				fmt.Printf("\n [no modified] %s", fileNames.Path)
				result.modTimes.Add(app.modTimes.Map[fileNames.Path])
				result.mapFile.Add(goformatter.MapFileRow{
					Names:   fileNames,
					ModTime: app.modTimes.Map[fileNames.Path].ModTime,
				})
				continue
			}
			changedFileChan <- fileNames
		}
	}
	return
}

func (app *App) processFile(changedFileChan chan names.FileNames, out chan processFileRow) {
	var (
		err      error
		modTimes = cache.NewModTimes()
		mapFile  = goformatter.NewMapFile("", "", "")
	)
	for {
		names, more := <-changedFileChan
		if !more {
			out <- processFileRow{
				modTimes: modTimes,
				mapFile:  mapFile,
			}
			return
		}
		var (
			bytes       []byte
			checksum    []byte
			lstat       os.FileInfo
			contentType string
		)
		firstLetter, _ := utf8.DecodeRuneInString(names.DirName)
		if !unicode.IsLetter(firstLetter) {
			panic(fmt.Sprintf("Package name must start from letter - your directory name is not start from letter (%s): %s", names.DirName, names.Path))
		}
		if bytes, err = app.FromFS.ReadFile(names.Path); err != nil {
			panic(err)
		}
		if checksum, err = calculateChecksum(bytes); err != nil {
			panic(err)
		}
		contentType = mime.TypeByExtension(filepath.Ext(names.Path))
		if contentType == "" {
			contentType = http.DetectContentType(bytes)
		}
		if lstat, err = app.FromFS.Lstat(names.Path); err != nil {
			panic(err)
		}
		file := &goformatter.BinaryFile{
			Names:       names,
			Bytes:       bytes,
			Checksum:    checksum,
			ModTime:     lstat.ModTime(),
			ContentType: contentType,
		}
		result := file.String()
		destPath := names.Path + ".go"
		if err = app.ToFS.WriteFile(destPath, []byte(result), filesystem.DefaultUnixFileMode); err != nil {
			panic(err)
		}
		modTimes.Add(cache.ModTime{
			Path:    names.Path,
			ModTime: file.ModTime,
		})
		mapFile.Add(goformatter.MapFileRow{
			Names:   names,
			ModTime: file.ModTime,
		})
		fmt.Printf("\n [generated] %s", destPath)
	}
}

func (app *App) Run() (err error) {
	var (
		fileChan       = make(chan names.FileNames, 2000)
		processed      = make(chan processFileRow, runtime.NumCPU())
		defaultDirName = names.ToUnderscore(filepath.Base(app.To))
		result         processFileRow
	)
	if defaultDirName == "" {
		defaultDirName = "fs"
	}
	result = processFileRow{
		modTimes: cache.NewModTimes(),
		mapFile:  goformatter.NewMapFile(app.PackagePrefix, defaultDirName, "FilesMap"),
	}
	if err = app.readModTimes(); err != nil {
		return
	}
	i := runtime.NumCPU()
	for ; i > 0; i-- {
		go app.processFile(fileChan, processed)
	}
	app.readDir(result, ".", fileChan)
	close(fileChan)

	for i := runtime.NumCPU(); i > 0; i-- {
		childResult := <-processed
		result.modTimes.Join(childResult.modTimes)
		result.mapFile.Join(childResult.mapFile.Rows)
	}
	close(processed)

	// math removed nodes
	removed := app.modTimes.Copy()
	removed.Except(result.modTimes)
	if len(removed.Map) > 0 {
		app.isRemoved = true
	}

	// remove all files created for removed nodes
	for path := range removed.Map {
		removedFilePath := path + ".go"
		app.ToFS.Remove(removedFilePath)
		fmt.Printf("\n [removed] %s", removedFilePath)
	}

	// save
	if app.isCreated || app.isModified || app.isRemoved {
		fmt.Printf("\n update mod times... ")
		if err = app.persistModTimes(result.modTimes); err != nil {
			return
		}
		fmt.Printf("ok")
		fmt.Printf("\n update map file... ")
		if err = app.persistMap(result.mapFile); err != nil {
			return
		}
		fmt.Printf("ok")
	}
	fmt.Printf("\n")
	return goaterr.ToError(app.errors)
}
