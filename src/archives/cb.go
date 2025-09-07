package archives

import (
	"path/filepath"
	"strings"

	"github.com/gen2brain/go-unarr"
)

func validImage(name string) bool {
	l := strings.ToLower(name)
	return strings.HasSuffix(l, ".jpg") || strings.HasSuffix(l, ".jpeg") || strings.HasSuffix(l, ".png") || strings.HasSuffix(l, ".webp")
}

func getBookInfoCB(path string) (BookInfo, error) {
	a, err := unarr.NewArchive(path)
	if err != nil {
		return BookInfo{}, err
	}
	defer a.Close()

	names, err := a.List()
	if err != nil {
		return BookInfo{}, err
	}

	pages := 0
	for _, name := range names {
		if strings.HasSuffix(name, "/") {
			continue
		}
		if validImage(name) {
			pages++
		}
	}

	title := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	return BookInfo{Title: title, Pages: pages}, nil
}
