package fsapp

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/SebastianPozoga/go-generator-filesystem/cache"
	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/filespace/memfs"
)

func TestCachedApp(t *testing.T) {
	var (
		fs, fromFS, toFS, cacheFS filesystem.Filespace
		resultBytes               []byte
		result                    string
		err                       error
	)
	if fs, err = memfs.NewFilespace(); err != nil {
		t.Error(err)
		return
	}
	if err = fs.MkdirAll("./from", filesystem.DefaultUnixDirMode); err != nil {
		t.Error(err)
		return
	}
	if err = fs.MkdirAll("./to", filesystem.DefaultUnixDirMode); err != nil {
		t.Error(err)
		return
	}
	if err = fs.WriteFile("./from/binaryfile.ex", []byte("12345"), filesystem.DefaultUnixFileMode); err != nil {
		t.Error(err)
		return
	}
	if fromFS, err = fs.Filespace("./from"); err != nil {
		t.Error(err)
		return
	}
	if toFS, err = fs.Filespace("./to"); err != nil {
		t.Error(err)
		return
	}
	if cacheFS, err = fs.Filespace("./cache"); err != nil {
		t.Error(err)
		return
	}
	var app = &App{
		From:    "./from",
		To:      "./to",
		Cache:   "./cache",
		FromFS:  fromFS,
		ToFS:    toFS,
		CacheFS: cacheFS,
	}
	if err = app.Run(); err != nil {
		t.Error(err)
		return
	}
	if !toFS.IsFile("binaryfile.ex.go") {
		t.Errorf("Expected file named binaryfile.ex.go into 'to' directory")
		return
	}
	if !cacheFS.IsFile(CacheModPath) {
		t.Errorf("Expected cached modtime")
		return
	}
	if resultBytes, err = toFS.ReadFile("binaryfile.ex.go"); err != nil {
		t.Error(err)
		return
	}
	result = string(resultBytes)
	if !strings.HasPrefix(result, "package to") {
		t.Errorf("Expected package named 'package to'")
		return
	}
	if !strings.Contains(result, "[]byte{49, 50, 51, 52, 53}") {
		t.Errorf("Expected file binaries []byte{49, 50, 51, 52, 53}")
		return
	}
}

func TestUpdateModTime(t *testing.T) {
	var (
		fs, fromFS, toFS, cacheFS filesystem.Filespace
		reader                    filesystem.Reader
		lstat                     os.FileInfo
		modTimeRow                cache.ModTime
		ok                        bool
		err                       error
	)
	if fs, err = memfs.NewFilespace(); err != nil {
		t.Error(err)
		return
	}
	if err = fs.MkdirAll("./from", filesystem.DefaultUnixDirMode); err != nil {
		t.Error(err)
		return
	}
	if err = fs.MkdirAll("./to", filesystem.DefaultUnixDirMode); err != nil {
		t.Error(err)
		return
	}
	if err = fs.WriteFile("./from/binaryfile.ex", []byte("12345"), filesystem.DefaultUnixFileMode); err != nil {
		t.Error(err)
		return
	}
	if fromFS, err = fs.Filespace("./from"); err != nil {
		t.Error(err)
		return
	}
	if toFS, err = fs.Filespace("./to"); err != nil {
		t.Error(err)
		return
	}
	if cacheFS, err = fs.Filespace("./cache"); err != nil {
		t.Error(err)
		return
	}
	if err = cacheFS.WriteFile(CacheModPath, []byte("binaryfile.ex;1991-01-01T00:00:00Z"), filesystem.DefaultUnixFileMode); err != nil {
		t.Error(err)
		return
	}
	var app = &App{
		From:    "./from",
		To:      "./to",
		Cache:   "./cache",
		FromFS:  fromFS,
		ToFS:    toFS,
		CacheFS: cacheFS,
	}
	if err = app.Run(); err != nil {
		t.Error(err)
		return
	}
	if !toFS.IsFile("binaryfile.ex.go") {
		t.Errorf("Expected file named binaryfile.ex.go into 'to' directory")
		return
	}
	if reader, err = cacheFS.Reader(CacheModPath); err != nil {
		reader.Close()
		t.Error(err)
		return
	}
	resultModTimes := cache.NewModTimes()
	if err = resultModTimes.Read(reader); err != nil {
		reader.Close()
		t.Error(err)
		return
	}
	reader.Close()
	if lstat, err = app.FromFS.Lstat("binaryfile.ex"); err != nil {
		t.Error(err)
		return
	}
	if modTimeRow, ok = resultModTimes.File("binaryfile.ex"); !ok {
		t.Errorf("Expected binaryfile.ex")
		return
	}
	if modTimeRow.ModTime.Unix() != lstat.ModTime().Unix() {
		var bytes []byte
		bytes, _ = cacheFS.ReadFile(CacheModPath)
		t.Errorf("Expected mod time from cache and mod time of file binaryfile.ex be equal. \n%s time from: %s\n modification time of file ./binaryfile.ex: %s\n %s file:\n%s",
			CacheModPath,
			modTimeRow.ModTime.Format(time.RFC3339),
			lstat.ModTime().Format(time.RFC3339),
			CacheModPath,
			string(bytes))
		return
	}
}
