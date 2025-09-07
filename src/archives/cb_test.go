package archives

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetBookInfoCBZ(t *testing.T) {
	path := filepath.Join("..", "..", "fixtures", "Full of Fun", "Full_of_Fun_001__Decker_Pub._1957.08__c2c___soothsayr_Yoc.cbz")
	book, err := getBookInfoCB(path)
	if err != nil {
		t.Fatalf("read cbz: %v", err)
	}
	if book.Pages <= 0 {
		t.Fatalf("pages = %d, want > 0", book.Pages)
	}
	if book.Title == "" {
		t.Fatalf("title is empty")
	}
	t.Logf("CBZ: Title=%s, Pages=%d", book.Title, book.Pages)
}

func TestGetBookInfoCBR(t *testing.T) {
	path := filepath.Join("..", "..", "fixtures", "Full of Fun", "Full_Of_Fun_001__c2c___1957___ABPC_.cbr")
	book, err := getBookInfoCB(path)
	if err != nil {
		t.Fatalf("read cbr: %v", err)
	}
	if book.Pages <= 0 {
		t.Fatalf("pages = %d, want > 0", book.Pages)
	}
	if book.Title == "" {
		t.Fatalf("title is empty")
	}
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
			if got := validImage(tt.filename); got != tt.want {
				t.Errorf("validImage(%q) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}

func TestExtractCBZ(t *testing.T) {
	// Setup
	inputPath := filepath.Join("..", "..", "fixtures", "Full of Fun", "Full_of_Fun_001__Decker_Pub._1957.08__c2c___soothsayr_Yoc.cbz")
	outputDir := t.TempDir()

	// Extract files
	extractedFiles, err := extractArchive(inputPath, outputDir)
	if err != nil {
		t.Fatalf("extractArchive() error = %v", err)
	}

	// Verify extraction results
	if len(extractedFiles) == 0 {
		t.Fatal("expected extracted files, got none")
	}

	// Check that all extracted files exist and are images
	for _, page := range extractedFiles {
		fullPath := filepath.Join(outputDir, page.Path)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("extracted file %s does not exist", page.Path)
		}

		// Verify it's an image file
		ext := strings.ToLower(filepath.Ext(page.Path))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
			t.Errorf("extracted file %s is not an image (ext: %s)", page.Path, ext)
		}

		// Verify dimensions are set
		if page.Width <= 0 || page.Height <= 0 {
			t.Errorf("invalid dimensions for %s: %dx%d", page.Path, page.Width, page.Height)
		}
	}

	// Verify files are sorted naturally
	expectedOrder := []string{"01.jpg", "02.jpg", "03.jpg", "04.jpg", "05.jpg", "06.jpg", "07.jpg", "08.jpg", "09.jpg", "10.jpg", "11.jpg", "12.jpg", "13.jpg", "14.jpg", "15.jpg", "16.jpg", "17.jpg", "18.jpg", "19.jpg", "20.jpg", "21.jpg", "22.jpg", "23.jpg", "24.jpg", "25.jpg", "26.jpg", "27.jpg", "28.jpg", "29.jpg", "30.jpg", "31.jpg", "32.jpg", "33.jpg", "34.jpg", "35.jpg", "36.jpg", "99-NFO.jpg"}

	if len(extractedFiles) != len(expectedOrder) {
		t.Errorf("expected %d files, got %d", len(expectedOrder), len(extractedFiles))
	}

	for i, page := range extractedFiles {
		if i < len(expectedOrder) && page.Path != expectedOrder[i] {
			t.Errorf("file order mismatch at index %d: got %s, want %s", i, page.Path, expectedOrder[i])
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
	if err != nil {
		t.Fatalf("extractArchive() error = %v", err)
	}

	// Verify extraction results
	if len(extractedFiles) == 0 {
		t.Fatal("expected extracted files, got none")
	}

	// Check that all extracted files exist and are images
	for _, page := range extractedFiles {
		fullPath := filepath.Join(outputDir, page.Path)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("extracted file %s does not exist", page.Path)
		}

		// Verify it's an image file
		ext := strings.ToLower(filepath.Ext(page.Path))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
			t.Errorf("extracted file %s is not an image (ext: %s)", page.Path, ext)
		}

		// Verify dimensions are set
		if page.Width <= 0 || page.Height <= 0 {
			t.Errorf("invalid dimensions for %s: %dx%d", page.Path, page.Width, page.Height)
		}
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
			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
