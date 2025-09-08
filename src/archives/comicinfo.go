package archives

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	"github.com/hekmon/go-comicinfo"
)

// parseComicInfo parses ComicInfo.xml content using the go-comicinfo library types
func parseComicInfo(xmlData []byte) (BookInfo, error) {
	// Try parsing as v2.1 first (most recent)
	var comicInfo comicinfo.ComicInfov21
	if err := xml.Unmarshal(xmlData, &comicInfo); err == nil {
		return convertComicInfoV21ToBookInfo(comicInfo), nil
	}

	// Try v2.0 format
	var comicInfoV2 comicinfo.ComicInfov2
	if err := xml.Unmarshal(xmlData, &comicInfoV2); err == nil {
		return convertComicInfoV2ToBookInfo(comicInfoV2), nil
	}

	// Try v1.0 format
	var comicInfoV1 comicinfo.ComicInfov1
	if err := xml.Unmarshal(xmlData, &comicInfoV1); err == nil {
		return convertComicInfoV1ToBookInfo(comicInfoV1), nil
	}

	return BookInfo{}, fmt.Errorf("failed to parse ComicInfo.xml")
}

// convertComicInfoV21ToBookInfo converts ComicInfov21 to BookInfo
func convertComicInfoV21ToBookInfo(ci comicinfo.ComicInfov21) BookInfo {
	bookInfo := BookInfo{}

	// Title - use Series and Number if Title is empty, or just Title
	if ci.Title != "" {
		bookInfo.Title = ci.Title
	} else if ci.Series != "" {
		if ci.Number > 0 {
			bookInfo.Title = fmt.Sprintf("%s #%d", ci.Series, ci.Number)
		} else {
			bookInfo.Title = ci.Series
		}
	}

	// Series information
	if ci.Series != "" {
		bookInfo.Series = ci.Series
	}

	if ci.Number > 0 {
		bookInfo.SeriesIndex = strconv.Itoa(ci.Number)
	}

	// Description
	if ci.Summary != "" {
		bookInfo.Description = ci.Summary
	}

	// Authors - combine creators
	var authors []string
	if ci.Writer != "" {
		authors = append(authors, splitCommaDelimited(ci.Writer)...)
	}
	if ci.Penciller != "" {
		authors = append(authors, splitCommaDelimited(ci.Penciller)...)
	}
	if ci.Inker != "" {
		authors = append(authors, splitCommaDelimited(ci.Inker)...)
	}
	if ci.Colorist != "" {
		authors = append(authors, splitCommaDelimited(ci.Colorist)...)
	}
	if ci.Letterer != "" {
		authors = append(authors, splitCommaDelimited(ci.Letterer)...)
	}
	if ci.Editor != "" {
		authors = append(authors, splitCommaDelimited(ci.Editor)...)
	}
	if ci.Translator != "" {
		authors = append(authors, splitCommaDelimited(ci.Translator)...)
	}
	// Remove duplicates
	authors = removeDuplicates(authors)
	if len(authors) > 0 {
		bookInfo.Authors = authors
	}

	// Publisher
	if ci.Publisher != "" {
		bookInfo.Publisher = ci.Publisher
	}

	// Published date - construct from Year, Month, Day
	if ci.Year > 0 {
		date := strconv.Itoa(ci.Year)
		if ci.Month > 0 {
			date += "-" + fmt.Sprintf("%02d", ci.Month)
			if ci.Day > 0 {
				date += "-" + fmt.Sprintf("%02d", ci.Day)
			}
		}
		bookInfo.PublishedDate = date
	}

	// Language
	if ci.LanguageISO != "" {
		bookInfo.Language = []string{ci.LanguageISO}
	}

	// Keywords - combine Genre, Tags, Characters, Teams, Locations
	var keywords []string
	if ci.Genre != "" {
		keywords = append(keywords, splitCommaDelimited(ci.Genre)...)
	}
	if ci.Tags != "" {
		keywords = append(keywords, splitCommaDelimited(ci.Tags)...)
	}
	if ci.Characters != "" {
		keywords = append(keywords, splitCommaDelimited(ci.Characters)...)
	}
	if ci.Teams != "" {
		keywords = append(keywords, splitCommaDelimited(ci.Teams)...)
	}
	if ci.Locations != "" {
		keywords = append(keywords, splitCommaDelimited(ci.Locations)...)
	}
	if len(keywords) > 0 {
		bookInfo.Keywords = removeDuplicates(keywords)
	}

	// Page count
	if ci.PageCount > 0 {
		bookInfo.Pages = ci.PageCount
	}

	return bookInfo
}

// convertComicInfoV2ToBookInfo converts ComicInfov2 to BookInfo
func convertComicInfoV2ToBookInfo(ci comicinfo.ComicInfov2) BookInfo {
	bookInfo := BookInfo{}

	// Title - use Series and Number if Title is empty, or just Title
	if ci.Title != "" {
		bookInfo.Title = ci.Title
	} else if ci.Series != "" {
		if ci.Number > 0 {
			bookInfo.Title = fmt.Sprintf("%s #%d", ci.Series, ci.Number)
		} else {
			bookInfo.Title = ci.Series
		}
	}

	// Series information
	if ci.Series != "" {
		bookInfo.Series = ci.Series
	}

	if ci.Number > 0 {
		bookInfo.SeriesIndex = strconv.Itoa(ci.Number)
	}

	// Description
	if ci.Summary != "" {
		bookInfo.Description = ci.Summary
	}

	// Authors - combine creators
	var authors []string
	if ci.Writer != "" {
		authors = append(authors, splitCommaDelimited(ci.Writer)...)
	}
	if ci.Penciller != "" {
		authors = append(authors, splitCommaDelimited(ci.Penciller)...)
	}
	if ci.Inker != "" {
		authors = append(authors, splitCommaDelimited(ci.Inker)...)
	}
	if ci.Colorist != "" {
		authors = append(authors, splitCommaDelimited(ci.Colorist)...)
	}
	if ci.Letterer != "" {
		authors = append(authors, splitCommaDelimited(ci.Letterer)...)
	}
	if ci.Editor != "" {
		authors = append(authors, splitCommaDelimited(ci.Editor)...)
	}
	// Remove duplicates
	authors = removeDuplicates(authors)
	if len(authors) > 0 {
		bookInfo.Authors = authors
	}

	// Publisher
	if ci.Publisher != "" {
		bookInfo.Publisher = ci.Publisher
	}

	// Published date - construct from Year, Month, Day
	if ci.Year > 0 {
		date := strconv.Itoa(ci.Year)
		if ci.Month > 0 {
			date += "-" + fmt.Sprintf("%02d", ci.Month)
			if ci.Day > 0 {
				date += "-" + fmt.Sprintf("%02d", ci.Day)
			}
		}
		bookInfo.PublishedDate = date
	}

	// Language
	if ci.LanguageISO != "" {
		bookInfo.Language = []string{ci.LanguageISO}
	}

	// Keywords - combine Genre, Characters, Teams, Locations
	var keywords []string
	if ci.Genre != "" {
		keywords = append(keywords, splitCommaDelimited(ci.Genre)...)
	}
	if ci.Characters != "" {
		keywords = append(keywords, splitCommaDelimited(ci.Characters)...)
	}
	if ci.Teams != "" {
		keywords = append(keywords, splitCommaDelimited(ci.Teams)...)
	}
	if ci.Locations != "" {
		keywords = append(keywords, splitCommaDelimited(ci.Locations)...)
	}
	if len(keywords) > 0 {
		bookInfo.Keywords = removeDuplicates(keywords)
	}

	// Page count
	if ci.PageCount > 0 {
		bookInfo.Pages = ci.PageCount
	}

	return bookInfo
}

// convertComicInfoV1ToBookInfo converts ComicInfov1 to BookInfo
func convertComicInfoV1ToBookInfo(ci comicinfo.ComicInfov1) BookInfo {
	bookInfo := BookInfo{}

	// Title - use Series and Number if Title is empty, or just Title
	if ci.Title != "" {
		bookInfo.Title = ci.Title
	} else if ci.Series != "" {
		if ci.Number > 0 {
			bookInfo.Title = fmt.Sprintf("%s #%d", ci.Series, ci.Number)
		} else {
			bookInfo.Title = ci.Series
		}
	}

	// Series information
	if ci.Series != "" {
		bookInfo.Series = ci.Series
	}

	if ci.Number > 0 {
		bookInfo.SeriesIndex = strconv.Itoa(ci.Number)
	}

	// Description
	if ci.Summary != "" {
		bookInfo.Description = ci.Summary
	}

	// Authors - combine creators
	var authors []string
	if ci.Writer != "" {
		authors = append(authors, splitCommaDelimited(ci.Writer)...)
	}
	if ci.Penciller != "" {
		authors = append(authors, splitCommaDelimited(ci.Penciller)...)
	}
	if ci.Inker != "" {
		authors = append(authors, splitCommaDelimited(ci.Inker)...)
	}
	if ci.Colorist != "" {
		authors = append(authors, splitCommaDelimited(ci.Colorist)...)
	}
	if ci.Letterer != "" {
		authors = append(authors, splitCommaDelimited(ci.Letterer)...)
	}
	if ci.Editor != "" {
		authors = append(authors, splitCommaDelimited(ci.Editor)...)
	}
	// Remove duplicates
	authors = removeDuplicates(authors)
	if len(authors) > 0 {
		bookInfo.Authors = authors
	}

	// Publisher
	if ci.Publisher != "" {
		bookInfo.Publisher = ci.Publisher
	}

	// Published date - construct from Year, Month
	if ci.Year > 0 {
		date := strconv.Itoa(ci.Year)
		if ci.Month > 0 {
			date += "-" + fmt.Sprintf("%02d", ci.Month)
		}
		bookInfo.PublishedDate = date
	}

	// Language
	if ci.Language != "" {
		bookInfo.Language = []string{ci.Language}
	}

	// Keywords - combine Genre
	var keywords []string
	if ci.Genre != "" {
		keywords = append(keywords, splitCommaDelimited(ci.Genre)...)
	}
	if len(keywords) > 0 {
		bookInfo.Keywords = removeDuplicates(keywords)
	}

	// Page count
	if ci.PageCount > 0 {
		bookInfo.Pages = ci.PageCount
	}

	return bookInfo
}

// splitCommaDelimited splits a comma-delimited string and trims whitespace
func splitCommaDelimited(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// removeDuplicates removes duplicate strings from a slice
func removeDuplicates(slice []string) []string {
	if len(slice) == 0 {
		return slice
	}

	seen := make(map[string]bool)
	result := make([]string, 0, len(slice))

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}
