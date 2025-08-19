package file

import (
	"bytes"
	"os"
	"strings"
)

// Combines all files content
func (f *FileLib) Merge(files []string) string {
	var sb strings.Builder
	for _, file := range files {
		if data, err := os.ReadFile(file); err == nil {
			sb.Write(data)
		}
	}
	return sb.String()
}

func (f *FileLib) MergeByte(files []string) []byte {
	var buf bytes.Buffer
	for _, file := range files {
		if data, err := os.ReadFile(file); err == nil {
			buf.Write(data)
		}
	}
	return buf.Bytes()
}
