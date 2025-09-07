// Package archives provides functions to handle book archives and PDFs
package archives

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	_ "golang.org/x/image/webp"
)

// Page represents a single page with its file path and dimensions
type Page struct {
	Path   string `json:"path"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

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
// Returns a list of extracted pages with file paths and dimensions
func Extract(inputFile, outputFolder string) ([]Page, error) {
	// Create output folder if it doesn't exist
	if err := os.MkdirAll(outputFolder, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output folder: %w", err)
	}

	// Determine file type and extract accordingly
	ext := strings.ToLower(filepath.Ext(inputFile))
	var extractedPages []Page
	var err error

	switch ext {
	case ".cbz", ".cbr", ".cb7", ".cbt":
		extractedPages, err = extractArchive(inputFile, outputFolder)
	case ".pdf":
		extractedPages, err = extractPDF(inputFile, outputFolder)
	default:
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}

	if err != nil {
		return nil, fmt.Errorf("extraction failed: %w", err)
	}

	return extractedPages, nil
}

// getImageDimensions returns the width and height of an image file
func getImageDimensions(imagePath string) (int, int, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	config, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}

	return config.Width, config.Height, nil
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
