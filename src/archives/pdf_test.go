package archives

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadMetadata(t *testing.T) {
	path := filepath.Join("..", "..", "fixtures", "testfile.pdf")

	book, err := getBookInfoPDF(path)
	require.NoError(t, err, "should successfully read PDF metadata")

	assert.Equal(t, 1, book.Pages, "should have exactly 1 page")
	assert.Equal(t, "Title of the Book", book.Title, "should have correct title")
	assert.Equal(t, []string{"The Author"}, book.Authors, "should have correct authors")
	assert.Equal(t, []string{"book", "fantasy"}, book.Keywords, "should have correct keywords")
}

func TestExtractPDF(t *testing.T) {
	// Setup
	inputPath := filepath.Join("..", "..", "fixtures", "testfile.pdf")
	outputDir := t.TempDir()

	// Extract files
	extractedFiles, err := extractPDF(inputPath, outputDir)
	require.NoError(t, err, "should successfully extract PDF")

	// Verify extraction results
	assert.NotEmpty(t, extractedFiles, "should extract files from PDF")

	// For the test PDF, we expect exactly 1 page
	expectedFiles := []string{"page_01.jpg"}
	assert.Len(t, extractedFiles, len(expectedFiles), "should extract expected number of files")

	// Check that all extracted files exist and are JPEG images
	for _, page := range extractedFiles {
		fullPath := filepath.Join(outputDir, page.Path)
		assert.FileExists(t, fullPath, "extracted file %s should exist", page.Path)

		// Verify it's a JPEG file
		ext := strings.ToLower(filepath.Ext(page.Path))
		assert.Contains(t, []string{".jpg", ".jpeg"}, ext,
			"extracted file %s should be a JPEG", page.Path)

		// Verify the file is not empty
		info, err := os.Stat(fullPath)
		require.NoError(t, err, "should be able to stat file %s", page.Path)
		assert.Greater(t, info.Size(), int64(0), "extracted file %s should not be empty", page.Path)

		// Verify dimensions are set
		assert.Greater(t, page.Width, 0, "width should be greater than 0 for %s", page.Path)
		assert.Greater(t, page.Height, 0, "height should be greater than 0 for %s", page.Path)
	}

	// Verify file naming convention (zero-padded)
	for i, page := range extractedFiles {
		if i+1 < 10 {
			assert.True(t, strings.HasPrefix(page.Path, "page_0"),
				"file %s should start with page_0 for single digits", page.Path)
		}
		expectedName := strings.Replace(page.Path, ".jpg", "", 1)
		if len(expectedName) > 1 {
			assert.True(t, strings.HasPrefix(expectedName, "page_0"),
				"file %s should be zero-padded", page.Path)
		}
	}

	t.Logf("Successfully extracted %d files from PDF", len(extractedFiles))
}

func TestExtractPDFErrorCases(t *testing.T) {
	tests := []struct {
		name        string
		inputPath   string
		outputDir   string
		expectError bool
	}{
		{
			name:        "nonexistent file",
			inputPath:   "nonexistent.pdf",
			outputDir:   t.TempDir(),
			expectError: true,
		},
		{
			name:        "invalid file format",
			inputPath:   "test.txt",
			outputDir:   t.TempDir(),
			expectError: true,
		},
		{
			name:        "nonexistent output directory",
			inputPath:   filepath.Join("..", "..", "fixtures", "testfile.pdf"),
			outputDir:   "/nonexistent/path/that/does/not/exist",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := extractPDF(tt.inputPath, tt.outputDir)
			if tt.expectError {
				assert.Error(t, err, "should return error for %s", tt.name)
			} else {
				assert.NoError(t, err, "should not return error for %s", tt.name)
			}
		})
	}
}

func TestExtractPDFWithUnidocLicense(t *testing.T) {
	// This test verifies that the PDF extraction works even without a license
	// (unidoc/unipdf has trial functionality)
	inputPath := filepath.Join("..", "..", "fixtures", "testfile.pdf")
	outputDir := t.TempDir()

	// Extract files
	extractedFiles, err := extractPDF(inputPath, outputDir)
	require.NoError(t, err, "should successfully extract PDF with unidoc/unipdf")

	// Should successfully extract at least one file
	assert.NotEmpty(t, extractedFiles, "should extract at least one file")

	// Verify the first file is a valid JPEG
	firstFile := filepath.Join(outputDir, extractedFiles[0].Path)
	info, err := os.Stat(firstFile)
	require.NoError(t, err, "should be able to stat first extracted file")

	// File should not be empty
	assert.Greater(t, info.Size(), int64(0), "extracted JPEG file should not be empty")

	// File should be reasonably sized (at least 1KB for a simple PDF page)
	assert.GreaterOrEqual(t, info.Size(), int64(1024),
		"extracted JPEG file should be at least 1KB, got %d bytes", info.Size())

	t.Logf("Successfully extracted %d files with unidoc/unipdf", len(extractedFiles))
}
