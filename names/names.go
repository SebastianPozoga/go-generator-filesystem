package names

import (
	"crypto/md5"
	"encoding/hex"
	"path/filepath"
	"strings"
)

type FileNames struct {
	Path        string
	FileName    string
	FileNameUCC string
	DirPath     string
	DirName     string
	DirNameU    string
	VarName     string
	HashDirPath string
}

func NewFileNames(path string, defaultDirName string) FileNames {
	path = filepath.Clean(path)
	path = strings.ReplaceAll(path, "\\", "/")
	var (
		filename    = filepath.Base(path)
		filenameUCC = ToCamelCase(filename, true)
		dirPath     = strings.ReplaceAll(filepath.Dir(path), "\\", "/")
		dirName     = filepath.Base(dirPath)
		dirNameU    string
		hashDirPath = hash(dirPath)
	)
	if dirName == "" || dirName == "." {
		dirName = defaultDirName
	}
	dirNameU = ToUnderscore(dirName)
	return FileNames{
		Path:        path,
		FileName:    filename,
		FileNameUCC: filenameUCC,
		DirPath:     dirPath,
		DirName:     dirName,
		DirNameU:    dirNameU,
		VarName:     "File" + filenameUCC,
		HashDirPath: hashDirPath,
	}
}

// hash create a new MD5 hash
func hash(input string) string {
	hash := md5.New()
	hash.Write([]byte(input))
	hashBytes := hash.Sum(nil)
	return "h" + hex.EncodeToString(hashBytes)
}
