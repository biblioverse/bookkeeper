package utils

import (
	"path/filepath"
	"strings"
)

func IsCBZ(path string) bool {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))
	return ext == "cbz"
}

func IsPDF(path string) bool {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))
	return ext == "pdf"
}

func IsValidBookFile(path string) bool {
	return IsCBZ(path) || IsPDF(path)
}
