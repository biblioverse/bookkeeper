// Package commands implements the bookkeeper commands
package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/biblioteca/bookkeeper/src/archives"
)

// PagesJSON represents the structure of the pages.json file
type PagesJSON struct {
	Pages []archives.Page `json:"pages"`
}

// Extract extracts files from an archive or PDF into the output folder
func Extract(inputFile, outputFolder string) error {
	// Use the archives package to extract files
	extractedPages, err := archives.Extract(inputFile, outputFolder)
	if err != nil {
		return fmt.Errorf("extraction failed: %w", err)
	}

	// Create pages.json
	if err := createPagesJSON(extractedPages, outputFolder); err != nil {
		return fmt.Errorf("failed to create pages.json: %w", err)
	}

	fmt.Printf("Extraction complete. %d files extracted to %s\n", len(extractedPages), outputFolder)
	return nil
}

// createPagesJSON creates the pages.json file with extracted pages and their dimensions
func createPagesJSON(pages []archives.Page, outputFolder string) error {
	pagesPath := filepath.Join(outputFolder, "pages.json")

	file, err := os.Create(pagesPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	pagesJSON := PagesJSON{Pages: pages}
	return encoder.Encode(pagesJSON)
}
