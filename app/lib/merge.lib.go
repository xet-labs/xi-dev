package lib

import (
	"bytes"
	"os"
	"strings"
)

type MergeLib struct{}

var Merge = &MergeLib{}

// Combines all files content
func (m *MergeLib) Files(files []string) string {
	var sb strings.Builder
	for _, file := range files {
		if data, err := os.ReadFile(file); err == nil {
			sb.Write(data)
		}
	}
	return sb.String()
}
func (m *MergeLib) FilesByte(files []string) []byte {
	var buf bytes.Buffer
	for _, file := range files {
		if data, err := os.ReadFile(file); err == nil {
			buf.Write(data)
		}
	}
	return buf.Bytes()
}
