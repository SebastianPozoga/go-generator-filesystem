package main

import (
	"fmt"
	"strings"
	"time"
)

type GOBinaryFile struct {
	Bytes       []byte
	Checksum    []byte
	Package     string
	Name        string
	ContentType string
	ModTime     time.Time
}

func (f *GOBinaryFile) String() string {
	var (
		builder strings.Builder
	)
	builder.WriteString("package " + f.Package + "\n\nimport \"time\"\n\nvar File" + f.Name + " = struct{")
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

func byteArray(builder *strings.Builder, bytes []byte) {
	builder.WriteString("[]byte{")
	for i := 0; i < len(bytes); i++ {
		builder.WriteString(fmt.Sprintf("%d", bytes[i]))
		if i < len(bytes)-1 {
			builder.WriteString(", ")
		}
	}
	builder.WriteString("}")
}

func formatTime(t time.Time) string {
	return fmt.Sprintf("time.Unix(%d, 0)", t.Unix())
}
