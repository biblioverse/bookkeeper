package archives

import (
	"fmt"
	"path/filepath"
	"strings"
)

type BookInfo struct {
	Title         string   `json:"title"`
	Pages         int      `json:"pages"`
	Authors       []string `json:"authors,omitempty"`
	Publisher     string   `json:"publisher,omitempty"`
	PublishedDate string   `json:"published_date,omitempty"`
	Keywords      []string `json:"keywords,omitempty"`
}

func GetBookInfo(path string) (BookInfo, error) {
	ext := filepath.Ext(path)
	switch strings.ToLower(ext) {
	case ".cbz", ".cbr":
		return getBookInfoCB(path)
	case ".pdf":
		return getBookInfoPDF(path)
	default:
		return BookInfo{}, fmt.Errorf("We don't know how to open this archive '%s'", path)
	}
}
