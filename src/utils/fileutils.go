package utils

import (
	"path/filepath"
	"strings"
)

// getFileExtension extracts and normalizes the file extension from a path
func getFileExtension(path string) string {
	return strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))
}

func IsValidBookFile(path string) bool {
	ext := getFileExtension(path)
	return ext == "cbz" || ext == "cbr" || ext == "cb7" || ext == "cbt" || ext == "pdf"
}
