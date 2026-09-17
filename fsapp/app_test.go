package fsapp

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
	if !strings.Contains(result, "//go:embed binaryfile.ex") {
		t.Errorf("Expected //go:embed directive for binaryfile.ex, got:\n%s", result)
		return
	}
	if !toFS.IsFile("binaryfile.ex") {
		t.Errorf("Expected raw asset binaryfile.ex to be written alongside the generated wrapper")
		return
	}
	if resultBytes, err = toFS.ReadFile("binaryfile.ex"); err != nil {
		t.Error(err)
		return
	}
	if string(resultBytes) != "12345" {
		t.Errorf("Expected raw asset content to be untouched, got: %q", string(resultBytes))
		return
	}
}

func TestAppIgnoresDirsAndFiles(t *testing.T) {
	var (
		fs, fromFS, toFS filesystem.Filespace
		mapBytes         []byte
		err              error
	)
	if fs, err = memfs.NewFilespace(); err != nil {
		t.Error(err)
		return
	}
	if err = fs.MkdirAll("./from/node_modules", filesystem.DefaultUnixDirMode); err != nil {
		t.Error(err)
		return
	}
	if err = fs.MkdirAll("./from/dist", filesystem.DefaultUnixDirMode); err != nil {
		t.Error(err)
		return
	}
	if err = fs.MkdirAll("./from/assets", filesystem.DefaultUnixDirMode); err != nil {
		t.Error(err)
		return
	}
	if err = fs.MkdirAll("./to", filesystem.DefaultUnixDirMode); err != nil {
		t.Error(err)
		return
	}
	if err = fs.WriteFile("./from/assets/keep.txt", []byte("keep"), filesystem.DefaultUnixFileMode); err != nil {
		t.Error(err)
		return
	}
	if err = fs.WriteFile("./from/node_modules/skip.txt", []byte("skip dir"), filesystem.DefaultUnixFileMode); err != nil {
		t.Error(err)
		return
	}
	if err = fs.WriteFile("./from/dist/build.txt", []byte("skip dir"), filesystem.DefaultUnixFileMode); err != nil {
		t.Error(err)
		return
	}
	if err = fs.WriteFile("./from/assets/secret.txt", []byte("skip file"), filesystem.DefaultUnixFileMode); err != nil {
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
		From:         "./from",
		To:           "./to",
		FromFS:       fromFS,
		ToFS:         toFS,
		IgnoredDirs:  []string{"node_modules", "dist"},
		IgnoredFiles: []string{"assets/secret.txt"},
	}
	if err = app.Run(); err != nil {
		t.Error(err)
		return
	}
	if !toFS.IsFile("assets/keep.txt.go") {
		t.Errorf("Expected assets/keep.txt.go to be generated")
		return
	}
	if toFS.IsFile("node_modules/skip.txt.go") {
		t.Errorf("Expected ignored node_modules file to be skipped")
		return
	}
	if toFS.IsFile("dist/build.txt.go") {
		t.Errorf("Expected ignored dist file to be skipped")
		return
	}
	if toFS.IsFile("assets/secret.txt.go") {
		t.Errorf("Expected ignored file to be skipped")
		return
	}
	if mapBytes, err = toFS.ReadFile("main.go"); err != nil {
		t.Error(err)
		return
	}
	mapContent := string(mapBytes)
	if !strings.Contains(mapContent, "assets/keep.txt") {
		t.Errorf("Expected map to include non-ignored file")
		return
	}
	if strings.Contains(mapContent, "node_modules/skip.txt") {
		t.Errorf("Expected map to exclude ignored directory files")
		return
	}
	if strings.Contains(mapContent, "assets/secret.txt") {
		t.Errorf("Expected map to exclude ignored files")
		return
	}
}

func TestAppGeneratesHiddenFiles(t *testing.T) {
	var (
		fs, fromFS, toFS filesystem.Filespace
		resultBytes      []byte
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
	if err = fs.WriteFile("./from/.gitignore", []byte("node_modules\n"), filesystem.DefaultUnixFileMode); err != nil {
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
	if !toFS.IsFile("gitignore.go") {
		t.Errorf("Expected hidden file to be generated as gitignore.go")
		return
	}
	if !toFS.IsFile(".gitignore") {
		t.Errorf("Expected raw asset to be written at its original hidden path .gitignore")
		return
	}
	if resultBytes, err = toFS.ReadFile("main.go"); err != nil {
		t.Error(err)
		return
	}
	result := string(resultBytes)
	if !strings.Contains(result, "\".gitignore\": FileGitignore") {
		t.Errorf("Expected map to reference FileGitignore for .gitignore, got:\n%s", result)
		return
	}
}

func TestAppGeneratesFilesFromHiddenDirectories(t *testing.T) {
	var (
		fs, fromFS, toFS filesystem.Filespace
		resultBytes      []byte
		err              error
	)
	if fs, err = memfs.NewFilespace(); err != nil {
		t.Error(err)
		return
	}
	if err = fs.MkdirAll("./from/web/node_modules/.bin", filesystem.DefaultUnixDirMode); err != nil {
		t.Error(err)
		return
	}
	if err = fs.MkdirAll("./to", filesystem.DefaultUnixDirMode); err != nil {
		t.Error(err)
		return
	}
	if err = fs.WriteFile("./from/web/node_modules/.bin/browserslist", []byte("12345"), filesystem.DefaultUnixFileMode); err != nil {
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
	if !toFS.IsFile("web/node_modules/.bin/browserslist.go") {
		t.Errorf("Expected hidden directory file to be generated")
		return
	}
	if !toFS.IsFile("web/node_modules/.bin/browserslist") {
		t.Errorf("Expected raw asset to be written at its original hidden-directory path")
		return
	}
	if resultBytes, err = toFS.ReadFile("web/node_modules/.bin/browserslist.go"); err != nil {
		t.Error(err)
		return
	}
	result := string(resultBytes)
	if !strings.HasPrefix(result, "package bin") {
		t.Errorf("Expected package named 'package bin', got:\n%s", result)
		return
	}
}
