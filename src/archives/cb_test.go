package archives

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBookInfoCBZ(t *testing.T) {
	path := filepath.Join("..", "..", "fixtures", "Full of Fun", "Full_of_Fun_001__Decker_Pub._1957.08__c2c___soothsayr_Yoc.cbz")

	book, err := getBookInfoCB(path)
	require.NoError(t, err, "should successfully read CBZ file")

	assert.Greater(t, book.Pages, 0, "should have more than 0 pages")
	assert.NotEmpty(t, book.Title, "title should not be empty")

	t.Logf("CBZ: Title=%s, Pages=%d", book.Title, book.Pages)
}

func TestGetBookInfoCBR(t *testing.T) {
	path := filepath.Join("..", "..", "fixtures", "Full of Fun", "Full_Of_Fun_001__c2c___1957___ABPC_.cbr")

	book, err := getBookInfoCB(path)
	require.NoError(t, err, "should successfully read CBR file")

	assert.Greater(t, book.Pages, 0, "should have more than 0 pages")
	assert.NotEmpty(t, book.Title, "title should not be empty")

	t.Logf("CBR: Title=%s, Pages=%d", book.Title, book.Pages)
}

func TestValidImage(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{"JPG", "image.jpg", true},
		{"JPEG", "image.jpeg", true},
		{"PNG", "image.png", true},
		{"WEBP", "image.webp", true},
		{"JPG uppercase", "IMAGE.JPG", true},
		{"JPEG uppercase", "IMAGE.JPEG", true},
		{"PNG uppercase", "IMAGE.PNG", true},
		{"WEBP uppercase", "IMAGE.WEBP", true},
		{"TXT", "readme.txt", false},
		{"Directory", "folder/", false},
		{"Empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validImage(tt.filename)
			assert.Equal(t, tt.want, got, "validImage(%q) should return %v", tt.filename, tt.want)
		})
	}
}

func TestExtractCBZ(t *testing.T) {
	// Setup
	inputPath := filepath.Join("..", "..", "fixtures", "Full of Fun", "Full_of_Fun_001__Decker_Pub._1957.08__c2c___soothsayr_Yoc.cbz")
	outputDir := t.TempDir()

	// Extract files
	extractedFiles, err := extractArchive(inputPath, outputDir)
	require.NoError(t, err, "should successfully extract CBZ archive")

	// Verify extraction results
	assert.NotEmpty(t, extractedFiles, "should extract files from archive")

	// Check that all extracted files exist and are images
	for _, page := range extractedFiles {
		fullPath := filepath.Join(outputDir, page.Path)
		assert.FileExists(t, fullPath, "extracted file %s should exist", page.Path)

		// Verify it's an image file
		ext := strings.ToLower(filepath.Ext(page.Path))
		assert.Contains(t, []string{".jpg", ".jpeg", ".png", ".webp"}, ext,
			"extracted file %s should be an image file", page.Path)

		// Verify dimensions are set
		assert.Greater(t, page.Width, 0, "width should be greater than 0 for %s", page.Path)
		assert.Greater(t, page.Height, 0, "height should be greater than 0 for %s", page.Path)
	}

	// Verify files are sorted naturally
	expectedOrder := []string{"01.jpg", "02.jpg", "03.jpg", "04.jpg", "05.jpg", "06.jpg", "07.jpg", "08.jpg", "09.jpg", "10.jpg", "11.jpg", "12.jpg", "13.jpg", "14.jpg", "15.jpg", "16.jpg", "17.jpg", "18.jpg", "19.jpg", "20.jpg", "21.jpg", "22.jpg", "23.jpg", "24.jpg", "25.jpg", "26.jpg", "27.jpg", "28.jpg", "29.jpg", "30.jpg", "31.jpg", "32.jpg", "33.jpg", "34.jpg", "35.jpg", "36.jpg", "99-NFO.jpg"}

	assert.Len(t, extractedFiles, len(expectedOrder), "should extract expected number of files")

	for i, page := range extractedFiles {
		if i < len(expectedOrder) {
			assert.Equal(t, expectedOrder[i], page.Path,
				"file order should match expected order at index %d", i)
		}
	}

	t.Logf("Successfully extracted %d files from CBZ", len(extractedFiles))
}

func TestExtractCBR(t *testing.T) {
	// Setup
	inputPath := filepath.Join("..", "..", "fixtures", "Full of Fun", "Full_Of_Fun_001__c2c___1957___ABPC_.cbr")
	outputDir := t.TempDir()

	// Extract files
	extractedFiles, err := extractArchive(inputPath, outputDir)
	require.NoError(t, err, "should successfully extract CBR archive")

	// Verify extraction results
	assert.NotEmpty(t, extractedFiles, "should extract files from archive")

	// Check that all extracted files exist and are images
	for _, page := range extractedFiles {
		fullPath := filepath.Join(outputDir, page.Path)
		assert.FileExists(t, fullPath, "extracted file %s should exist", page.Path)

		// Verify it's an image file
		ext := strings.ToLower(filepath.Ext(page.Path))
		assert.Contains(t, []string{".jpg", ".jpeg", ".png", ".webp"}, ext,
			"extracted file %s should be an image file", page.Path)

		// Verify dimensions are set
		assert.Greater(t, page.Width, 0, "width should be greater than 0 for %s", page.Path)
		assert.Greater(t, page.Height, 0, "height should be greater than 0 for %s", page.Path)
	}

	t.Logf("Successfully extracted %d files from CBR", len(extractedFiles))
}

func TestExtractArchiveErrorCases(t *testing.T) {
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
			_, err := extractArchive(tt.inputPath, tt.outputDir)
			if tt.expectError {
				assert.Error(t, err, "should return error for %s", tt.name)
			} else {
				assert.NoError(t, err, "should not return error for %s", tt.name)
			}
		})
	}
}
