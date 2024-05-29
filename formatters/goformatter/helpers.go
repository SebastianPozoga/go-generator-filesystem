package goformatter

import (
	"fmt"
	"io"
	"time"
)

func byteArray(w io.Writer, bytes []byte) {
	io.WriteString(w, "[]byte{")
	for i := 0; i < len(bytes); i++ {
		io.WriteString(w, fmt.Sprintf("%d", bytes[i]))
		if i < len(bytes)-1 {
			io.WriteString(w, ", ")
		}
	}
	io.WriteString(w, "}")
}

func formatTime(t time.Time) string {
	return fmt.Sprintf("time.Unix(%d, 0)", t.Unix())
}
