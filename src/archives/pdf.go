package archives

import (
	"path/filepath"
	"strings"

	"rsc.io/pdf"
)

func getBookInfoPDF(path string) (BookInfo, error) {
	r, err := pdf.Open(path)
	if err != nil {
		return BookInfo{}, err
	}

	pages := r.NumPage()
	title := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	var authors []string
	var keywords []string

	if info := r.Trailer().Key("Info"); info.Kind() == pdf.Dict {
		if v := strings.TrimSpace(info.Key("Author").Text()); v != "" {
			authors = []string{v}
		}
		if v := strings.TrimSpace(info.Key("Title").Text()); v != "" && v != ".pdf" {
			title = v
		}
		if v := strings.TrimSpace(info.Key("Keywords").Text()); v != "" {
			for _, s := range strings.Split(v, ",") {
				s = strings.TrimSpace(s)
				if s != "" {
					keywords = append(keywords, s)
				}
			}
		}
	}

	return BookInfo{
		Title:    title,
		Pages:    pages,
		Authors:  authors,
		Keywords: keywords,
	}, nil
}
