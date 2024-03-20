package main

import (
	"bytes"
	"fmt"
)

type GOBinaryFile struct {
	Bytes   []byte
	Package string
	Name    string
}

func (f *GOBinaryFile) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("package " + f.Package + "\n\nvar " + f.Name + " []byte = []byte{")
	for i := 0; i < len(f.Bytes); i++ {
		buffer.WriteString(fmt.Sprintf("%d", f.Bytes[i]))
		if i < len(f.Bytes)-1 {
			buffer.WriteString(", ")
		}
	}
	buffer.WriteString("}\n")
	return buffer.String()
}
