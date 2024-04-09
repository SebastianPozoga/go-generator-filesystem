package main

import (
	"strings"
	"testing"

	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/filespace/memfs"
)

func TestApp(t *testing.T) {
	var (
		fs, fromFS, toFS filesystem.Filespace
		resultBytes      []byte
		result           string
		err              error
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
	var app = &App{
		From:   "./from",
		To:     "./to",
		FromFS: fromFS,
		ToFS:   toFS,
	}
	if err = app.Run(); err != nil {
		t.Error(err)
		return
	}
	if !toFS.IsFile("binaryfile.ex.go") {
		t.Errorf("Expected file named binaryfile.ex.go into 'to' directory")
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
