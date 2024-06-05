package to

import (
	"time"
)
var Checksum []byte = []byte{212, 189, 10, 168, 3, 76, 128, 10, 236, 191, 173, 43, 220, 55, 123, 61, 200, 29, 213, 63, 65, 1, 104, 153, 54, 232, 167, 188, 15, 2, 160, 102}

type FilesMapType struct {
	Checksum []byte
	Data []byte
	ContentType string
	ModTime time.Time
}

var FilesMap = map[string]FilesMapType {
	"from.txt": FileFromTxt,
}
