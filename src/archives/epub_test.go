package archives

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBookInfoEPUB(t *testing.T) {
	path := filepath.Join("..", "..", "fixtures", "pg11-images-3.epub")

	book, err := getBookInfoEPUB(path)
	require.NoError(t, err, "should successfully read EPUB file")

	assert.Greater(t, book.Pages, 0, "should have more than 0 pages")
	assert.Equal(t, "Alice's Adventures in Wonderland", book.Title, "title should not be empty")
	assert.Equal(t, []string{"Lewis Carroll"}, book.Authors, "should have correct author")
	assert.Equal(t, "2008-06-27", book.PublishedDate, "should have correct publication date")
	assert.Equal(t, []string{
		"Fantasy fiction",
		"Children's stories",
		"Imaginary places -- Juvenile fiction",
		"Alice (Fictitious character from Carroll) -- Juvenile fiction",
	}, book.Keywords, "should have correct keywords")

	t.Logf("EPUB: Title=%s, Pages=%d, Authors=%v, PublishedDate=%s, Keywords=%v",
		book.Title, book.Pages, book.Authors, book.PublishedDate, book.Keywords)
}

func TestGetBookInfoEPUB2(t *testing.T) {
	path := filepath.Join("..", "..", "fixtures", "pg76832-images.epub")

	book, err := getBookInfoEPUB(path)
	require.NoError(t, err, "should successfully read EPUB file")

	assert.Greater(t, book.Pages, 0, "should have more than 0 pages")
	assert.Equal(t, "Sämtliche Werke 21: Der Spieler. Der ewige Gatte.", book.Title, "title should not be empty")
	assert.Equal(t, []string{"Fyodor Dostoyevsky"}, book.Authors, "should have correct author")
	assert.Equal(t, "2025-09-07", book.PublishedDate, "should have correct publication date")
	assert.Equal(t, []string(nil), book.Keywords, "should have correct keywords")

	t.Logf("EPUB: Title=%s, Pages=%d, Authors=%v, PublishedDate=%s",
		book.Title, book.Pages, book.Authors, book.PublishedDate)
}
