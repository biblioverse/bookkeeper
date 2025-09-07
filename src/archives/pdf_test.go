package archives

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadMetadata(t *testing.T) {
	path := filepath.Join("..", "..", "fixtures", "testfile.pdf")
	book, err := getBookInfoPDF(path)
	if err != nil {
		t.Fatalf("read pdf: %v", err)
	}
	if book.Pages != 1 {
		t.Fatalf("pages = %d, want 1", book.Pages)
	}
	if book.Title != "Title of the Book" {
		t.Fatalf("title = %q, want %q", book.Title, "Title of the Book")
	}
	if len(book.Authors) != 1 || book.Authors[0] != "The Author" {
		t.Fatalf("authors = %#v, want [The Author]", book.Authors)
	}
	if len(book.Keywords) != 2 || book.Keywords[0] != "book" || book.Keywords[1] != "fantasy" {
		t.Fatalf("keywords = %#v, want [book fantasy]", book.Keywords)
	}
}

func TestExtractPDF(t *testing.T) {
	// Setup
	inputPath := filepath.Join("..", "..", "fixtures", "testfile.pdf")
	outputDir := t.TempDir()

	// Extract files
	extractedFiles, err := extractPDF(inputPath, outputDir)
	if err != nil {
		t.Fatalf("extractPDF() error = %v", err)
	}

	// Verify extraction results
	if len(extractedFiles) == 0 {
		t.Fatal("expected extracted files, got none")
	}

	// For the test PDF, we expect exactly 1 page
	expectedFiles := []string{"page_01.jpg"}
	if len(extractedFiles) != len(expectedFiles) {
		t.Errorf("expected %d files, got %d", len(expectedFiles), len(extractedFiles))
	}

	// Check that all extracted files exist and are JPEG images
	for _, file := range extractedFiles {
		fullPath := filepath.Join(outputDir, file)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("extracted file %s does not exist", file)
		}

		// Verify it's a JPEG file
		ext := strings.ToLower(filepath.Ext(file))
		if ext != ".jpg" && ext != ".jpeg" {
			t.Errorf("extracted file %s is not a JPEG (ext: %s)", file, ext)
		}

		// Verify the file is not empty
		info, err := os.Stat(fullPath)
		if err != nil {
			t.Errorf("failed to stat file %s: %v", file, err)
		}
		if info.Size() == 0 {
			t.Errorf("extracted file %s is empty", file)
		}
	}

	// Verify file naming convention (zero-padded)
	for i, file := range extractedFiles {
		expectedName := strings.Replace(file, ".jpg", "", 1)
		if !strings.HasPrefix(expectedName, "page_0") && len(expectedName) > 1 {
			t.Errorf("file %s should be zero-padded", file)
		}
		if i+1 < 10 && !strings.HasPrefix(file, "page_0") {
			t.Errorf("file %s should start with page_0 for single digits", file)
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
			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
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
	if err != nil {
		t.Fatalf("extractPDF() error = %v", err)
	}

	// Should successfully extract at least one file
	if len(extractedFiles) == 0 {
		t.Fatal("expected at least one extracted file")
	}

	// Verify the first file is a valid JPEG
	firstFile := filepath.Join(outputDir, extractedFiles[0])
	info, err := os.Stat(firstFile)
	if err != nil {
		t.Fatalf("failed to stat first extracted file: %v", err)
	}

	// File should not be empty
	if info.Size() == 0 {
		t.Error("extracted JPEG file is empty")
	}

	// File should be reasonably sized (at least 1KB for a simple PDF page)
	if info.Size() < 1024 {
		t.Errorf("extracted JPEG file is too small: %d bytes", info.Size())
	}

	t.Logf("Successfully extracted %d files with unidoc/unipdf", len(extractedFiles))
}
