package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/filespace/diskfs"
	"github.com/goatcms/goatcore/filesystem/fsloop"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// App is main application object. Contains all app data
type App struct {
	From, To         string
	FS, FromFS, ToFS filesystem.Filespace
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
}

func (app *App) Run() error {
	loop := fsloop.NewLoop(&fsloop.LoopData{
		Filespace: app.FromFS,
		DirFilter: func(fs filesystem.Filespace, subPath string) bool {
			return strings.HasSuffix(subPath, ".ex")
		},
		OnDir: func(fs filesystem.Filespace, subPath string) error {
			return fs.MkdirAll(subPath, filesystem.DefaultUnixDirMode)
		},
		OnFile: func(fs filesystem.Filespace, subPath string) error {
			var (
				err      error
				bytes    []byte
				filename = toCamelCase(filepath.Base(subPath), true)
				dirname  = toUnderscore(filepath.Base(filepath.Dir(subPath)))
			)
			if bytes, err = fs.ReadFile(subPath); err != nil {
				panic(err)
			}
			if dirname == "" {
				dirname = toUnderscore(filepath.Base(app.To))
			}
			if dirname == "" {
				dirname = "fs"
			}
			file := &GOBinaryFile{
				Bytes:   bytes,
				Package: dirname,
				Name:    filename,
			}
			firstLetter, _ := utf8.DecodeRuneInString(file.Package)
			if !unicode.IsLetter(firstLetter) {
				panic(fmt.Sprintf("Package name must start from letter - your directory name is not start from letter (%s)", subPath))
			}
			firstLetter, _ = utf8.DecodeRuneInString(file.Name)
			if !unicode.IsLetter(firstLetter) {
				panic(fmt.Sprintf("Function name must start from letter - your file name is not start from letter (%s)", subPath))
			}
			result := file.String()
			return app.ToFS.WriteFile(subPath+".go", []byte(result), filesystem.DefaultUnixFileMode)
		},
		Consumers:  runtime.NumCPU(),
		Producents: 1,
	}, nil)
	loop.Run("./")
	loop.Wait()
	return goaterr.ToError(loop.Errors())
}
