package goformatter

import (
	"strings"
	"time"

	"github.com/SebastianPozoga/go-generator-filesystem/names"
)

type BinaryFile struct {
	Names       names.FileNames
	Bytes       []byte
	Checksum    []byte
	ContentType string
	ModTime     time.Time
}

func (f *BinaryFile) String() string {
	var (
		builder strings.Builder
		dataVar = f.Names.VarName + "Data"
	)
	builder.WriteString("package " + f.Names.DirNameU + "\n\nimport (\n\t_ \"embed\"\n\t\"time\"\n)\n\n")
	builder.WriteString("//go:embed " + f.Names.FileName + "\n")
	builder.WriteString("var " + dataVar + " []byte\n\n")
	builder.WriteString("var " + f.Names.VarName + " = struct{")
	builder.WriteString("\n\tChecksum []byte")
	builder.WriteString("\n\tData []byte")
	builder.WriteString("\n\tContentType string")
	builder.WriteString("\n\tModTime time.Time")
	builder.WriteString("\n}{\n\t")
	byteArray(&builder, []byte(f.Checksum))
	builder.WriteString(",\n\t")
	builder.WriteString(dataVar)
	builder.WriteString(",\n\t\"")
	builder.WriteString(f.ContentType)
	builder.WriteString("\",\n\t")
	builder.WriteString(formatTime(f.ModTime))
	builder.WriteString(",\n}")
	return builder.String()
}
