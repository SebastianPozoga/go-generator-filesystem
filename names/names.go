package names

import (
	"path/filepath"
)

type FileNames struct {
	Path        string
	FileName    string
	FileNameUCC string
	DirPath     string
	DirName     string
	DirNameU    string
	VarName     string
}

func NewFileNames(path string, defaultDirName string) FileNames {
	path = filepath.Clean(path)
	var (
		filename    = filepath.Base(path)
		filenameUCC = ToCamelCase(filename, true)
		dirPath     = filepath.Dir(path)
		dirName     = filepath.Base(dirPath)
		dirNameU    string
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
	}
}
