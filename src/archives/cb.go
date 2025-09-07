package archives

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gen2brain/go-unarr"
	"github.com/maruel/natural"
)

func validImage(name string) bool {
	l := strings.ToLower(name)
	return strings.HasSuffix(l, ".jpg") || strings.HasSuffix(l, ".jpeg") || strings.HasSuffix(l, ".png") || strings.HasSuffix(l, ".webp")
}

func getBookInfoCB(path string) (BookInfo, error) {
	a, err := unarr.NewArchive(path)
	if err != nil {
		return BookInfo{}, err
	}
	defer a.Close()

	names, err := a.List()
	if err != nil {
		return BookInfo{}, err
	}

	pages := 0
	for _, name := range names {
		if strings.HasSuffix(name, "/") {
			continue
		}
		if validImage(name) {
			pages++
		}
	}

	title := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	return BookInfo{Title: title, Pages: pages}, nil
}

// extractArchive extracts files from archive formats (CBZ, CBR, etc.)
func extractArchive(inputFile, outputFolder string) ([]string, error) {
	archive, err := unarr.NewArchive(inputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open archive: %w", err)
	}
	defer archive.Close()

	// Extract all files to the output folder
	extractedFiles, err := archive.Extract(outputFolder)
	if err != nil {
		return nil, fmt.Errorf("failed to extract archive: %w", err)
	}

	// Filter out directories and convert to relative paths
	var relativeFiles []string
	for _, filePath := range extractedFiles {
		// Skip directories
		if strings.HasSuffix(filePath, "/") {
			continue
		}

		// The Extract method returns relative paths, so we can use them directly
		// But let's make sure they're clean relative paths
		cleanPath := filepath.Clean(filePath)
		relativeFiles = append(relativeFiles, cleanPath)
	}

	// Apply natural sorting for archive files
	sort.Slice(relativeFiles, func(i, j int) bool {
		return natural.Less(relativeFiles[i], relativeFiles[j])
	})

	return relativeFiles, nil
}
