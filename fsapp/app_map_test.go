package fsapp

import (
	"strings"
	"testing"

	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/filespace/memfs"
)

func TestMapedApp(t *testing.T) {
	var (
		fs, fromFS, toFS, cacheFS filesystem.Filespace
		mapBytes                  []byte
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
	if err = fs.WriteFile("./from/subdir/binaryfile.ex", []byte("12345"), filesystem.DefaultUnixFileMode); err != nil {
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
		PackagePrefix: "github.com/user/repo/",
		From:          "./from",
		To:            "./to",
		Cache:         "./cache",
		FromFS:        fromFS,
		ToFS:          toFS,
		CacheFS:       cacheFS,
	}
	if err = app.Run(); err != nil {
		t.Error(err)
		return
	}
	if !toFS.IsFile("main.go") {
		t.Errorf("Expected file named main.go into 'to' directory")
		return
	}
	if mapBytes, err = app.ToFS.ReadFile("main.go"); err != nil {
		t.Error(err)
		return
	}
	mapContent := string(mapBytes)
	if !strings.HasPrefix(mapContent, "package to") {
		t.Errorf("Expected package named 'package to' and take:\n%s", mapContent)
		return
	}
	if !strings.Contains(mapContent, "binaryfile.ex") {
		t.Errorf("Expected map file contanins 'binaryfile.ex' and take:\n%s", mapContent)
		return
	}
	if !strings.Contains(mapContent, "github.com/user/repo/") {
		t.Errorf("Expected map file contanins 'binaryfile.ex' and take:\n%s", mapContent)
		return
	}

}
