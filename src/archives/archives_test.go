package archives

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBookInfoIntegration(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{"EPUB Alice", "pg11-images-3.epub", false},
		{"EPUB Dostoyevsky", "pg76832-images.epub", false},
		{"PDF", "testfile.pdf", false},
		{"CBZ", filepath.Join("Full of Fun", "Full_of_Fun_001__Decker_Pub._1957.08__c2c___soothsayr_Yoc.cbz"), false},
		{"CBR", filepath.Join("Full of Fun", "Full_Of_Fun_001__c2c___1957___ABPC_.cbr"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join("..", "..", "fixtures", tt.filename)

			book, err := GetBookInfo(path)
			if tt.wantErr {
				assert.Error(t, err, "should return error for %s", tt.name)
			} else {
				assert.NoError(t, err, "should not return error for %s", tt.name)
				assert.NotEmpty(t, book.Title, "title should not be empty for %s", tt.name)
				assert.GreaterOrEqual(t, book.Pages, 0, "pages should be 0 or more for %s", tt.name)

				t.Logf("%s: Title=%s, Pages=%d, Authors=%v",
					tt.name, book.Title, book.Pages, book.Authors)
			}
		})
	}
}

func TestIsValidBookFileExtended(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{"CBZ", "book.cbz", true},
		{"CBR", "book.cbr", true},
		{"CB7", "book.cb7", true},
		{"CBT", "book.cbt", true},
		{"PDF", "book.pdf", true},
		{"EPUB", "book.epub", true},
		{"EPUB uppercase", "BOOK.EPUB", true},
		{"TXT", "readme.txt", false},
		{"DOCX", "document.docx", false},
		{"Empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidBookFile(tt.filename)
			assert.Equal(t, tt.want, got, "IsValidBookFile(%q) should return %v", tt.filename, tt.want)
		})
	}
}
