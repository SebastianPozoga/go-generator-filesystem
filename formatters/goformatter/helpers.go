package goformatter

import (
	"fmt"
	"strings"
	"time"
)

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
