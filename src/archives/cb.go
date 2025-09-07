package archives

import (
	"archive/zip"
	"fmt"
	"path/filepath"
	"strings"
)

func validImage(name string) bool {
	l := strings.ToLower(name)
	return strings.HasSuffix(l, ".jpg") || strings.HasSuffix(l, ".jpeg") || strings.HasSuffix(l, ".png") || strings.HasSuffix(l, ".webp")
}

func getBookInfoCB(path string) (BookInfo, error) {
	if strings.ToLower(filepath.Ext(path)) != ".cbz" {
		return BookInfo{}, fmt.Errorf("unsupported comic archive: %s", path)
	}

	r, err := zip.OpenReader(path)
	if err != nil {
		return BookInfo{}, err
	}
	defer r.Close()

	pages := 0
	for _, f := range r.File {
		if !f.FileInfo().IsDir() && validImage(f.Name) {
			pages++
		}
	}

	title := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	return BookInfo{Title: title, Pages: pages}, nil
}
