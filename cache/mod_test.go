package cache

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestRead(t *testing.T) {
	var modTimes = NewModTimes()
	expectedFiles := []struct {
		Path string
		Time time.Time
	}{
		{"/path/to/file1", time.Date(2024, 11, 22, 0, 0, 0, 0, time.UTC)},
		{"/path/to/file2", time.Date(2024, 11, 23, 0, 0, 0, 0, time.UTC)},
		{"/path/to/file3", time.Date(2024, 12, 24, 0, 0, 0, 0, time.UTC)},
	}
	reader := strings.NewReader(`/path/to/file1;2024-11-22T00:00:00Z
/path/to/file2;2024-11-23T00:00:00Z
/path/to/file3;2024-12-24T00:00:00Z`)
	modTimes.Read(reader)
	for _, expected := range expectedFiles {
		if _, ok := modTimes.Map[expected.Path]; !ok {
			t.Errorf("Expected file %s", expected.Path)
		}
		if !expected.Time.Equal(modTimes.Map[expected.Path].ModTime) {
			t.Errorf("Expected time %s and take %s", expected.Time.Format(time.RFC3339), modTimes.Map[expected.Path].ModTime.Format(time.RFC3339))
		}
	}
}

func TestWriter(t *testing.T) {
	var modTimes = &ModTimes{
		Map: map[string]ModTime{
			"/path/to/file1": {"/path/to/file1", time.Date(2024, 11, 22, 0, 0, 0, 0, time.UTC)},
			"/path/to/file2": {"/path/to/file2", time.Date(2024, 11, 23, 0, 0, 0, 0, time.UTC)},
			"/path/to/file3": {"/path/to/file3", time.Date(2024, 12, 24, 0, 0, 0, 0, time.UTC)},
		},
	}
	expecteds := []string{
		"/path/to/file1;2024-11-22T00:00:00Z",
		"/path/to/file2;2024-11-23T00:00:00Z",
		"/path/to/file3;2024-12-24T00:00:00Z",
	}
	buf := bytes.NewBufferString("")
	modTimes.Write(buf)
	result := buf.String()
	for _, expected := range expecteds {
		if !strings.Contains(result, expected) {
			t.Errorf("Expected %s in \n%s", expected, result)
			return
		}
	}
}
