package commands

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExtractCBZ(t *testing.T) {
	// Setup
	inputPath := filepath.Join("..", "..", "fixtures", "Full of Fun", "Full_of_Fun_001__Decker_Pub._1957.08__c2c___soothsayr_Yoc.cbz")
	outputDir := t.TempDir()

	// Extract files using the main Extract function
	err := Extract(inputPath, outputDir)
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	// Verify pages.json was created
	pagesPath := filepath.Join(outputDir, "pages.json")
	if _, err := os.Stat(pagesPath); os.IsNotExist(err) {
		t.Fatal("pages.json was not created")
	}

	// Read and verify pages.json content
	pagesData, err := os.ReadFile(pagesPath)
	if err != nil {
		t.Fatalf("failed to read pages.json: %v", err)
	}

	var pages []string
	if err := json.Unmarshal(pagesData, &pages); err != nil {
		t.Fatalf("failed to unmarshal pages.json: %v", err)
	}

	if len(pages) == 0 {
		t.Fatal("pages.json is empty")
	}

	// Verify all pages are image files
	for _, page := range pages {
		ext := strings.ToLower(filepath.Ext(page))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
			t.Errorf("page %s is not an image file", page)
		}

		// Verify the file exists
		fullPath := filepath.Join(outputDir, page)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("page file %s does not exist", page)
		}
	}

	t.Logf("Successfully extracted %d pages from CBZ", len(pages))
}

func TestExtractCBR(t *testing.T) {
	// Setup
	inputPath := filepath.Join("..", "..", "fixtures", "Full of Fun", "Full_Of_Fun_001__c2c___1957___ABPC_.cbr")
	outputDir := t.TempDir()

	// Extract files using the main Extract function
	err := Extract(inputPath, outputDir)
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	// Verify pages.json was created
	pagesPath := filepath.Join(outputDir, "pages.json")
	if _, err := os.Stat(pagesPath); os.IsNotExist(err) {
		t.Fatal("pages.json was not created")
	}

	// Read and verify pages.json content
	pagesData, err := os.ReadFile(pagesPath)
	if err != nil {
		t.Fatalf("failed to read pages.json: %v", err)
	}

	var pages []string
	if err := json.Unmarshal(pagesData, &pages); err != nil {
		t.Fatalf("failed to unmarshal pages.json: %v", err)
	}

	if len(pages) == 0 {
		t.Fatal("pages.json is empty")
	}

	// Verify all pages are image files
	for _, page := range pages {
		ext := strings.ToLower(filepath.Ext(page))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
			t.Errorf("page %s is not an image file", page)
		}

		// Verify the file exists
		fullPath := filepath.Join(outputDir, page)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("page file %s does not exist", page)
		}
	}

	t.Logf("Successfully extracted %d pages from CBR", len(pages))
}

func TestExtractPDF(t *testing.T) {
	// Setup
	inputPath := filepath.Join("..", "..", "fixtures", "testfile.pdf")
	outputDir := t.TempDir()

	// Extract files using the main Extract function
	err := Extract(inputPath, outputDir)
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	// Verify pages.json was created
	pagesPath := filepath.Join(outputDir, "pages.json")
	if _, err := os.Stat(pagesPath); os.IsNotExist(err) {
		t.Fatal("pages.json was not created")
	}

	// Read and verify pages.json content
	pagesData, err := os.ReadFile(pagesPath)
	if err != nil {
		t.Fatalf("failed to read pages.json: %v", err)
	}

	var pages []string
	if err := json.Unmarshal(pagesData, &pages); err != nil {
		t.Fatalf("failed to unmarshal pages.json: %v", err)
	}

	// For the test PDF, we expect exactly 1 page
	if len(pages) != 1 {
		t.Errorf("expected 1 page, got %d", len(pages))
	}

	// Verify the page is a JPEG file
	page := pages[0]
	ext := strings.ToLower(filepath.Ext(page))
	if ext != ".jpg" && ext != ".jpeg" {
		t.Errorf("page %s is not a JPEG file (ext: %s)", page, ext)
	}

	// Verify the file exists and is not empty
	fullPath := filepath.Join(outputDir, page)
	info, err := os.Stat(fullPath)
	if err != nil {
		t.Errorf("page file %s does not exist: %v", page, err)
	}
	if info.Size() == 0 {
		t.Errorf("page file %s is empty", page)
	}

	// Verify file naming convention (zero-padded)
	if !strings.HasPrefix(page, "page_0") {
		t.Errorf("page file %s should be zero-padded", page)
	}

	t.Logf("Successfully extracted %d pages from PDF", len(pages))
}

func TestExtractErrorCases(t *testing.T) {
	tests := []struct {
		name        string
		inputPath   string
		outputDir   string
		expectError bool
	}{
		{
			name:        "nonexistent file",
			inputPath:   "nonexistent.cbz",
			outputDir:   t.TempDir(),
			expectError: true,
		},
		{
			name:        "unsupported file format",
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
			err := Extract(tt.inputPath, tt.outputDir)
			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestCreatePagesJSON(t *testing.T) {
	// Setup
	outputDir := t.TempDir()
	files := []string{"01.jpg", "02.jpg", "03.jpg"}

	// Create pages.json
	err := createPagesJSON(files, outputDir)
	if err != nil {
		t.Fatalf("createPagesJSON() error = %v", err)
	}

	// Verify pages.json was created
	pagesPath := filepath.Join(outputDir, "pages.json")
	if _, err := os.Stat(pagesPath); os.IsNotExist(err) {
		t.Fatal("pages.json was not created")
	}

	// Read and verify content
	pagesData, err := os.ReadFile(pagesPath)
	if err != nil {
		t.Fatalf("failed to read pages.json: %v", err)
	}

	var result []string
	if err := json.Unmarshal(pagesData, &result); err != nil {
		t.Fatalf("failed to unmarshal pages.json: %v", err)
	}

	// Verify content matches
	if len(result) != len(files) {
		t.Errorf("expected %d files, got %d", len(files), len(result))
	}

	for i, file := range files {
		if i < len(result) && result[i] != file {
			t.Errorf("file at index %d: got %s, want %s", i, result[i], file)
		}
	}
}

func TestCreatePagesJSONErrorCases(t *testing.T) {
	tests := []struct {
		name      string
		files     []string
		outputDir string
		expectErr bool
	}{
		{
			name:      "nonexistent directory",
			files:     []string{"01.jpg"},
			outputDir: "/nonexistent/path/that/does/not/exist",
			expectErr: true,
		},
		{
			name:      "empty files list",
			files:     []string{},
			outputDir: t.TempDir(),
			expectErr: false,
		},
		{
			name:      "nil files list",
			files:     nil,
			outputDir: t.TempDir(),
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := createPagesJSON(tt.files, tt.outputDir)
			if tt.expectErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
