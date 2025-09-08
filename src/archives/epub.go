package archives

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pirmd/epub"
)

// getBookInfoEPUB extracts metadata from EPUB files
func getBookInfoEPUB(path string) (BookInfo, error) {
	// Use the simpler Information method
	info, err := epub.GetMetadataFromFile(path)
	if err != nil {
		return BookInfo{}, fmt.Errorf("failed to read EPUB metadata: %w", err)
	}

	// Extract title - use filename as fallback
	title := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	if len(info.Title) > 0 && len(info.Title[0]) > 0 {
		title = string(info.Title[0])
	}

	subtitle := []string{}
	if len(info.SubTitle) > 0 {
		subtitle = info.SubTitle
	}

	language := []string{}
	if len(info.Language) > 0 {
		language = info.Language
	}

	description := ""
	if len(info.Description) > 0 {
		description = strings.Join(info.Description, ", ")
	}

	var series string
	var seriesIndex string
	if len(info.Series) > 0 {
		series = string(info.Series)
		if len(info.SeriesIndex) > 0 {
			seriesIndex = info.SeriesIndex
		}
	}

	// Extract authors
	var authors []string
	for _, creator := range info.Creator {
		if creator.FullName != "" {
			authors = append(authors, creator.FullName)
		}
	}

	// Extract publisher
	publisher := ""
	if len(info.Publisher) > 0 && len(info.Publisher[0]) > 0 {
		publisher = string(info.Publisher[0])
	}

	// Extract publication date
	publishedDate := ""
	if len(info.Date) > 0 {
		publishedDate = info.Date[0].Stamp
	}

	// Extract subjects/keywords
	var keywords []string
	for _, subject := range info.Subject {
		if len(subject) > 0 {
			keywords = append(keywords, string(subject))
		}
	}

	// For EPUB, we don't have a direct page count - spine items give us an approximation
	// We need to open the package to get spine information
	book, err := epub.Open(path)
	if err != nil {
		return BookInfo{}, fmt.Errorf("failed to open EPUB file: %w", err)
	}
	defer book.Close()

	pkg, err := book.Package()
	if err != nil {
		return BookInfo{}, fmt.Errorf("failed to read EPUB package: %w", err)
	}

	pages := len(pkg.Spine.Itemrefs)

	return BookInfo{
		Title:         title,
		SubTitle:      subtitle,
		Language:      language,
		Description:   description,
		Series:        series,
		SeriesIndex:   seriesIndex,
		Pages:         pages,
		Authors:       authors,
		Publisher:     publisher,
		PublishedDate: publishedDate,
		Keywords:      keywords,
	}, nil
}
