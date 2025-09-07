// Package commands implements the bookkeeper commands
package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/biblioteca/bookkeeper/src/archives"
)

// Extract extracts files from an archive or PDF into the output folder
func Extract(inputFile, outputFolder string) error {
	// Use the archives package to extract files
	extractedFiles, err := archives.Extract(inputFile, outputFolder)
	if err != nil {
		return fmt.Errorf("extraction failed: %w", err)
	}

	// Create pages.json
	if err := createPagesJSON(extractedFiles, outputFolder); err != nil {
		return fmt.Errorf("failed to create pages.json: %w", err)
	}

	fmt.Printf("Extraction complete. %d files extracted to %s\n", len(extractedFiles), outputFolder)
	return nil
}

// createPagesJSON creates the pages.json file with extracted file paths
func createPagesJSON(files []string, outputFolder string) error {
	pagesPath := filepath.Join(outputFolder, "pages.json")

	file, err := os.Create(pagesPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(files)
}
