package commands

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/biblioteca/bookkeeper/src/archives"
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

	var pagesJSON PagesJSON
	if err := json.Unmarshal(pagesData, &pagesJSON); err != nil {
		t.Fatalf("failed to unmarshal pages.json: %v", err)
	}

	if len(pagesJSON.Pages) == 0 {
		t.Fatal("pages.json contains no pages")
	}

	// Verify all pages are image files
	for _, page := range pagesJSON.Pages {
		ext := strings.ToLower(filepath.Ext(page.Path))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
			t.Errorf("page %s is not an image file", page.Path)
		}

		// Verify the file exists
		fullPath := filepath.Join(outputDir, page.Path)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("page file %s does not exist", page.Path)
		}

		// Verify dimensions are set
		if page.Width <= 0 || page.Height <= 0 {
			t.Errorf("invalid dimensions for %s: %dx%d", page.Path, page.Width, page.Height)
		}
	}

	t.Logf("Successfully extracted %d pages from CBZ", len(pagesJSON.Pages))
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

	var pagesJSON PagesJSON
	if err := json.Unmarshal(pagesData, &pagesJSON); err != nil {
		t.Fatalf("failed to unmarshal pages.json: %v", err)
	}

	if len(pagesJSON.Pages) == 0 {
		t.Fatal("pages.json contains no pages")
	}

	// Verify all pages are image files
	for _, page := range pagesJSON.Pages {
		ext := strings.ToLower(filepath.Ext(page.Path))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
			t.Errorf("page %s is not an image file", page.Path)
		}

		// Verify the file exists
		fullPath := filepath.Join(outputDir, page.Path)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("page file %s does not exist", page.Path)
		}

		// Verify dimensions are set
		if page.Width <= 0 || page.Height <= 0 {
			t.Errorf("invalid dimensions for %s: %dx%d", page.Path, page.Width, page.Height)
		}
	}

	t.Logf("Successfully extracted %d pages from CBR", len(pagesJSON.Pages))
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

	var pagesJSON PagesJSON
	if err := json.Unmarshal(pagesData, &pagesJSON); err != nil {
		t.Fatalf("failed to unmarshal pages.json: %v", err)
	}

	// For the test PDF, we expect exactly 1 page
	if len(pagesJSON.Pages) != 1 {
		t.Errorf("expected 1 page, got %d", len(pagesJSON.Pages))
	}

	// Verify the page is a JPEG file
	page := pagesJSON.Pages[0]
	ext := strings.ToLower(filepath.Ext(page.Path))
	if ext != ".jpg" && ext != ".jpeg" {
		t.Errorf("page %s is not a JPEG file (ext: %s)", page.Path, ext)
	}

	// Verify the file exists and is not empty
	fullPath := filepath.Join(outputDir, page.Path)
	info, err := os.Stat(fullPath)
	if err != nil {
		t.Errorf("page file %s does not exist: %v", page.Path, err)
	}
	if info.Size() == 0 {
		t.Errorf("page file %s is empty", page.Path)
	}

	// Verify file naming convention (zero-padded)
	if !strings.HasPrefix(page.Path, "page_0") {
		t.Errorf("page file %s should be zero-padded", page.Path)
	}

	// Verify dimensions are set
	if page.Width <= 0 || page.Height <= 0 {
		t.Errorf("invalid dimensions for %s: %dx%d", page.Path, page.Width, page.Height)
	}

	t.Logf("Successfully extracted %d pages from PDF", len(pagesJSON.Pages))
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
	pages := []archives.Page{
		{Path: "01.jpg", Width: 800, Height: 600},
		{Path: "02.jpg", Width: 1024, Height: 768},
		{Path: "03.jpg", Width: 1200, Height: 900},
	}

	// Create pages.json
	err := createPagesJSON(pages, outputDir)
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

	var result PagesJSON
	if err := json.Unmarshal(pagesData, &result); err != nil {
		t.Fatalf("failed to unmarshal pages.json: %v", err)
	}

	// Verify content matches
	if len(result.Pages) != len(pages) {
		t.Errorf("expected %d pages, got %d", len(pages), len(result.Pages))
	}

	for i, expectedPage := range pages {
		if i < len(result.Pages) {
			resultPage := result.Pages[i]
			if resultPage.Path != expectedPage.Path {
				t.Errorf("page at index %d: got path %s, want %s", i, resultPage.Path, expectedPage.Path)
			}
			if resultPage.Width != expectedPage.Width {
				t.Errorf("page at index %d: got width %d, want %d", i, resultPage.Width, expectedPage.Width)
			}
			if resultPage.Height != expectedPage.Height {
				t.Errorf("page at index %d: got height %d, want %d", i, resultPage.Height, expectedPage.Height)
			}
		}
	}
}

func TestCreatePagesJSONErrorCases(t *testing.T) {
	tests := []struct {
		name      string
		pages     []archives.Page
		outputDir string
		expectErr bool
	}{
		{
			name:      "nonexistent directory",
			pages:     []archives.Page{{Path: "01.jpg", Width: 800, Height: 600}},
			outputDir: "/nonexistent/path/that/does/not/exist",
			expectErr: true,
		},
		{
			name:      "empty pages list",
			pages:     []archives.Page{},
			outputDir: t.TempDir(),
			expectErr: false,
		},
		{
			name:      "nil pages list",
			pages:     nil,
			outputDir: t.TempDir(),
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := createPagesJSON(tt.pages, tt.outputDir)
			if tt.expectErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
