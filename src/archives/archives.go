// Package archives provides functions to handle book archives and PDFs
package archives

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// BookInfo holds metadata about a book
type BookInfo struct {
	Title         string   `json:"title"`
	Pages         int      `json:"pages"`
	Authors       []string `json:"authors,omitempty"`
	Publisher     string   `json:"publisher,omitempty"`
	PublishedDate string   `json:"published_date,omitempty"`
	Keywords      []string `json:"keywords,omitempty"`
}

// GetBookInfo retrieves metadata from a book archive or PDF file
func GetBookInfo(path string) (BookInfo, error) {
	ext := filepath.Ext(path)
	switch strings.ToLower(ext) {
	case ".cbz", ".cbr", ".cb7", ".cbt":
		return getBookInfoCB(path)
	case ".pdf":
		return getBookInfoPDF(path)
	default:
		return BookInfo{}, fmt.Errorf("we don't know how to open this archive '%s'", path)
	}
}

// Extract extracts files from an archive or PDF into the output folder
// Returns a list of extracted file paths relative to the output folder
func Extract(inputFile, outputFolder string) ([]string, error) {
	// Create output folder if it doesn't exist
	if err := os.MkdirAll(outputFolder, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output folder: %w", err)
	}

	// Determine file type and extract accordingly
	ext := strings.ToLower(filepath.Ext(inputFile))
	var extractedFiles []string
	var err error

	switch ext {
	case ".cbz", ".cbr", ".cb7", ".cbt":
		extractedFiles, err = extractArchive(inputFile, outputFolder)
	case ".pdf":
		extractedFiles, err = extractPDF(inputFile, outputFolder)
	default:
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}

	if err != nil {
		return nil, fmt.Errorf("extraction failed: %w", err)
	}

	return extractedFiles, nil
}

// getFileExtension extracts and normalizes the file extension from a path
func getFileExtension(path string) string {
	return strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))
}

// IsValidBookFile checks if the file has a valid book file extension
func IsValidBookFile(path string) bool {
	ext := getFileExtension(path)
	return ext == "cbz" || ext == "cbr" || ext == "cb7" || ext == "cbt" || ext == "pdf"
}
