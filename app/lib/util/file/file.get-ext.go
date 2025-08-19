package file

import (
	"os"
	"path/filepath"
)

// ListFilesWithExt returns a slice of file paths with the given extension from baseDir recursively.
// Example: ListFilesWithExt("static/css", ".css")
func (f *FileLib) GetWithExt(ext string, baseDirs ...string) ([]string, error) {
	var files []string

	for _, baseDir := range baseDirs {
		if err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() && filepath.Ext(info.Name()) == ext {
				files = append(files, path)
			}
			return nil
		}); err != nil {
			return nil, err
		}
	}

	return files, nil
}

func (f *FileLib) GetWithExts(exts []string, baseDirs ...string) ([]string, error) {
	var files []string
	extMap := make(map[string]bool, len(exts))

	for _, ext := range exts {
		if ext != "" {
			extMap[ext] = true
		}
	}

	for _, baseDir := range baseDirs {
		if err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
			if err != nil && !info.IsDir() && extMap[filepath.Ext(info.Name())] {
				files = append(files, path)
			}
			return nil
		}); err != nil {
			return nil, err
		}
	}

	return files, nil
}
