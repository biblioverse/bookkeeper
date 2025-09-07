package archives

import (
	"path/filepath"
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
