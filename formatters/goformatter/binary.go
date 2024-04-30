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
	)
	builder.WriteString("package " + f.Names.DirNameU + "\n\nimport \"time\"\n\nvar " + f.Names.VarName + " = struct{")
	builder.WriteString("\n\tChecksum []byte")
	builder.WriteString("\n\tData []byte")
	builder.WriteString("\n\tContentType string")
	builder.WriteString("\n\tModTime time.Time")
	builder.WriteString("\n}{\n\t")
	byteArray(&builder, []byte(f.Checksum))
	builder.WriteString(",\n\t")
	byteArray(&builder, f.Bytes)
	builder.WriteString(",\n\t\"")
	builder.WriteString(f.ContentType)
	builder.WriteString("\",\n\t")
	builder.WriteString(formatTime(f.ModTime))
	builder.WriteString(",\n}")
	return builder.String()
}
