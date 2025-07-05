package util

import (
	"os"
	"path/filepath"
)

// ListFilesWithExt returns a slice of file paths with the given extension from baseDir recursively.
// Example: ListFilesWithExt("static/css", ".css")
func GetFilesWithExt(ext string, baseDirs ...string) ([]string, error) {
	var files []string

	for _, baseDir := range baseDirs {
		err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() && filepath.Ext(info.Name()) == ext {
				files = append(files, path)
			}
			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	return files, nil
}

