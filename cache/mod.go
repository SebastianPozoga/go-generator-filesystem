package cache

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/goatcms/goatcore/filesystem"
)

const ModTimeSeparator = ";"

type ModTime struct {
	Path    string
	ModTime time.Time
}

type ModTimes struct {
	Map map[string]ModTime
}

func NewModTimes() *ModTimes {
	return &ModTimes{
		Map: make(map[string]ModTime),
	}
}

func (md *ModTimes) Read(file io.Reader) (err error) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var modificationTime time.Time
		line := scanner.Text()
		parts := strings.Split(line, ModTimeSeparator)
		if len(parts) != 2 {
			return fmt.Errorf("incorrect line formatt: %v", line)
		}
		filePath := parts[0]
		modificationTimeStr := parts[1]
		if modificationTime, err = time.Parse(time.RFC3339, modificationTimeStr); err != nil {
			return fmt.Errorf("parse time error: %v", err)
		}
		row := ModTime{
			Path:    filePath,
			ModTime: modificationTime,
		}
		md.Map[row.Path] = row
	}
	return
}

func (md *ModTimes) Write(file io.Writer) (err error) {
	i := 0
	for _, row := range md.Map {
		if i != 0 {
			file.Write([]byte("\n"))
		}
		i++
		timeStr := row.ModTime.Format(time.RFC3339)
		file.Write([]byte(row.Path))
		file.Write([]byte(ModTimeSeparator))
		file.Write([]byte(timeStr))
	}
	return
}

func (md *ModTimes) Add(row ModTime) {
	md.Map[row.Path] = row
}

func (md *ModTimes) Remove(path string) {
	delete(md.Map, path)
}

func (md *ModTimes) Copy() (copied *ModTimes) {
	copied = &ModTimes{
		Map: make(map[string]ModTime, len(md.Map)),
	}
	for k, v := range md.Map {
		copied.Map[k] = v
	}
	return
}

func (md *ModTimes) Join(row *ModTimes) {
	for k, v := range row.Map {
		md.Map[k] = v
	}
}

func (md *ModTimes) Except(second *ModTimes) {
	for k := range second.Map {
		delete(md.Map, k)
	}
}

func (md *ModTimes) File(path string) (mt ModTime, ok bool) {
	mt, ok = md.Map[path]
	return
}

func (md *ModTimes) IsFileModified(fs filesystem.Filespace, subPath string) (modified, created bool) {
	fileMT, ok := md.File(subPath)
	if !ok {
		return false, true
	}
	info, err := fs.Lstat(subPath)
	if err != nil {
		return true, false
	}
	return info.ModTime().Unix() != fileMT.ModTime.Unix(), false
}
