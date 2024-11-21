package goformatter

import (
	"crypto/sha256"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/SebastianPozoga/go-generator-filesystem/names"
)

type MapFileRow struct {
	Names   names.FileNames
	ModTime time.Time
}

type MapFile struct {
	Package       string
	PackagePrefix string
	VarName       string
	Rows          map[string]MapFileRow
}

func NewMapFile(prefix, packageName, varName string) *MapFile {
	return &MapFile{
		PackagePrefix: prefix,
		Package:       packageName,
		VarName:       varName,
		Rows:          make(map[string]MapFileRow),
	}
}

func (f *MapFile) Join(ex map[string]MapFileRow) {
	for k, v := range ex {
		f.Rows[k] = v
	}
}

func (f *MapFile) Add(row MapFileRow) {
	f.Rows[row.Names.Path] = row
}

func (f *MapFile) checksum() (sum []byte, err error) {
	var keys = make([]string, 0, len(f.Rows))
	for key := range f.Rows {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	hash := sha256.New()
	for _, key := range keys {
		var bytes []byte
		if bytes, err = f.Rows[key].ModTime.MarshalBinary(); err != nil {
			return
		}
		if _, err = hash.Write(bytes); err != nil {
			return
		}
	}
	return hash.Sum(nil), nil
}

func (f *MapFile) Write(w io.Writer) {
	var (
		checksum []byte
		err      error
	)
	if checksum, err = f.checksum(); err != nil {
		panic(err)
	}
	io.WriteString(w, "package "+f.Package)
	packages := map[string]string{}
	for _, v := range f.Rows {
		if v.Names.DirPath == "." {
			continue
		}
		packages[v.Names.DirPath] = v.Names.HashDirPath
	}
	io.WriteString(w, "\n\nimport (\n\t\"time\"")
	for path, name := range packages {
		io.WriteString(w, "\n\t"+name)
		io.WriteString(w, " \""+f.PackagePrefix)
		io.WriteString(w, path)
		io.WriteString(w, "\"")
	}
	io.WriteString(w, "\n)")
	io.WriteString(w, "\nvar Checksum []byte = ")
	byteArray(w, checksum)
	io.WriteString(w, "\n\ntype "+f.VarName+"Type struct {")
	io.WriteString(w, "\n\tChecksum []byte")
	io.WriteString(w, "\n\tData []byte")
	io.WriteString(w, "\n\tContentType string")
	io.WriteString(w, "\n\tModTime time.Time")
	io.WriteString(w, "\n}")
	io.WriteString(w, "\n\nvar "+f.VarName+" = map[string]"+f.VarName+"Type {")
	for _, row := range f.Rows {
		if row.Names.DirPath == "." {
			io.WriteString(w, fmt.Sprintf("\n\t\"%s\": %s,", row.Names.Path, row.Names.VarName))
			continue
		}
		io.WriteString(w, fmt.Sprintf("\n\t\"%s\": %s.%s,", row.Names.Path, row.Names.HashDirPath, row.Names.VarName))
	}
	io.WriteString(w, "\n}\n")
}

func (f *MapFile) String() string {
	var builder strings.Builder
	f.Write(&builder)
	return builder.String()
}
