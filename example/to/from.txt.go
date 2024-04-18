package to

import "time"

var FromTxt = struct{
	Checksum []byte
	Data []byte
	ContentType string
	ModTime time.Time
}{
	[]byte{242, 155, 198, 74, 157, 55, 50, 180, 185, 3, 81, 37, 253, 179, 40, 95, 91, 100, 85, 119, 142, 220, 167, 36, 20, 103, 30, 12, 163, 178, 224, 222},
	[]byte{84, 104, 105, 115, 32, 105, 115, 32, 97, 32, 116, 101, 115, 116, 32, 102, 105, 108, 101, 46},
	"text/plain; charset=utf-8",
	time.Unix(1713424464, 0),
}