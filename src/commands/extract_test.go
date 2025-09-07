package commands

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/biblioteca/bookkeeper/src/archives"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractCBZ(t *testing.T) {
	// Setup
	inputPath := filepath.Join("..", "..", "fixtures", "Full of Fun", "Full_of_Fun_001__Decker_Pub._1957.08__c2c___soothsayr_Yoc.cbz")
	outputDir := t.TempDir()

	// Extract files using the main Extract function
	err := Extract(inputPath, outputDir)
	require.NoError(t, err, "should successfully extract CBZ file")

	// Verify pages.json was created
	pagesPath := filepath.Join(outputDir, "pages.json")
	assert.FileExists(t, pagesPath, "pages.json should be created")

	// Read and verify pages.json content
	pagesData, err := os.ReadFile(pagesPath)
	require.NoError(t, err, "should be able to read pages.json")

	var pagesJSON PagesJSON
	err = json.Unmarshal(pagesData, &pagesJSON)
	require.NoError(t, err, "should be able to unmarshal pages.json")

	assert.NotEmpty(t, pagesJSON.Pages, "pages.json should contain pages")

	// Verify all pages are image files
	for _, page := range pagesJSON.Pages {
		ext := strings.ToLower(filepath.Ext(page.Path))
		assert.Contains(t, []string{".jpg", ".jpeg", ".png", ".webp"}, ext,
			"page %s should be an image file", page.Path)

		// Verify the file exists
		fullPath := filepath.Join(outputDir, page.Path)
		assert.FileExists(t, fullPath, "page file %s should exist", page.Path)

		// Verify dimensions are set
		assert.Greater(t, page.Width, 0, "width should be greater than 0 for %s", page.Path)
		assert.Greater(t, page.Height, 0, "height should be greater than 0 for %s", page.Path)
	}

	t.Logf("Successfully extracted %d pages from CBZ", len(pagesJSON.Pages))
}

func TestExtractCBR(t *testing.T) {
	// Setup
	inputPath := filepath.Join("..", "..", "fixtures", "Full of Fun", "Full_Of_Fun_001__c2c___1957___ABPC_.cbr")
	outputDir := t.TempDir()

	// Extract files using the main Extract function
	err := Extract(inputPath, outputDir)
	require.NoError(t, err, "should successfully extract CBR file")

	// Verify pages.json was created
	pagesPath := filepath.Join(outputDir, "pages.json")
	assert.FileExists(t, pagesPath, "pages.json should be created")

	// Read and verify pages.json content
	pagesData, err := os.ReadFile(pagesPath)
	require.NoError(t, err, "should be able to read pages.json")

	var pagesJSON PagesJSON
	err = json.Unmarshal(pagesData, &pagesJSON)
	require.NoError(t, err, "should be able to unmarshal pages.json")

	assert.NotEmpty(t, pagesJSON.Pages, "pages.json should contain pages")

	// Verify all pages are image files
	for _, page := range pagesJSON.Pages {
		ext := strings.ToLower(filepath.Ext(page.Path))
		assert.Contains(t, []string{".jpg", ".jpeg", ".png", ".webp"}, ext,
			"page %s should be an image file", page.Path)

		// Verify the file exists
		fullPath := filepath.Join(outputDir, page.Path)
		assert.FileExists(t, fullPath, "page file %s should exist", page.Path)

		// Verify dimensions are set
		assert.Greater(t, page.Width, 0, "width should be greater than 0 for %s", page.Path)
		assert.Greater(t, page.Height, 0, "height should be greater than 0 for %s", page.Path)
	}

	t.Logf("Successfully extracted %d pages from CBR", len(pagesJSON.Pages))
}

func TestExtractPDF(t *testing.T) {
	// Setup
	inputPath := filepath.Join("..", "..", "fixtures", "testfile.pdf")
	outputDir := t.TempDir()

	// Extract files using the main Extract function
	err := Extract(inputPath, outputDir)
	require.NoError(t, err, "should successfully extract PDF file")

	// Verify pages.json was created
	pagesPath := filepath.Join(outputDir, "pages.json")
	assert.FileExists(t, pagesPath, "pages.json should be created")

	// Read and verify pages.json content
	pagesData, err := os.ReadFile(pagesPath)
	require.NoError(t, err, "should be able to read pages.json")

	var pagesJSON PagesJSON
	err = json.Unmarshal(pagesData, &pagesJSON)
	require.NoError(t, err, "should be able to unmarshal pages.json")

	// For the test PDF, we expect exactly 1 page
	assert.Len(t, pagesJSON.Pages, 1, "should extract exactly 1 page from test PDF")

	// Verify the page is a JPEG file
	page := pagesJSON.Pages[0]
	ext := strings.ToLower(filepath.Ext(page.Path))
	assert.Contains(t, []string{".jpg", ".jpeg"}, ext,
		"page %s should be a JPEG file", page.Path)

	// Verify the file exists and is not empty
	fullPath := filepath.Join(outputDir, page.Path)
	info, err := os.Stat(fullPath)
	require.NoError(t, err, "page file %s should exist", page.Path)
	assert.Greater(t, info.Size(), int64(0), "page file %s should not be empty", page.Path)

	// Verify file naming convention (zero-padded)
	assert.True(t, strings.HasPrefix(page.Path, "page_0"),
		"page file %s should be zero-padded", page.Path)

	// Verify dimensions are set
	assert.Greater(t, page.Width, 0, "width should be greater than 0 for %s", page.Path)
	assert.Greater(t, page.Height, 0, "height should be greater than 0 for %s", page.Path)

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
			if tt.expectError {
				assert.Error(t, err, "should return error for %s", tt.name)
			} else {
				assert.NoError(t, err, "should not return error for %s", tt.name)
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
	require.NoError(t, err, "should successfully create pages.json")

	// Verify pages.json was created
	pagesPath := filepath.Join(outputDir, "pages.json")
	assert.FileExists(t, pagesPath, "pages.json should be created")

	// Read and verify content
	pagesData, err := os.ReadFile(pagesPath)
	require.NoError(t, err, "should be able to read pages.json")

	var result PagesJSON
	err = json.Unmarshal(pagesData, &result)
	require.NoError(t, err, "should be able to unmarshal pages.json")

	// Verify content matches
	assert.Len(t, result.Pages, len(pages), "should have correct number of pages")

	for i, expectedPage := range pages {
		if i < len(result.Pages) {
			resultPage := result.Pages[i]
			assert.Equal(t, expectedPage.Path, resultPage.Path,
				"page path should match at index %d", i)
			assert.Equal(t, expectedPage.Width, resultPage.Width,
				"page width should match at index %d", i)
			assert.Equal(t, expectedPage.Height, resultPage.Height,
				"page height should match at index %d", i)
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
			if tt.expectErr {
				assert.Error(t, err, "should return error for %s", tt.name)
			} else {
				assert.NoError(t, err, "should not return error for %s", tt.name)
			}
		})
	}
}
