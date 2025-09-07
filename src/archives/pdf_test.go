package archives

import (
	"path/filepath"
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
