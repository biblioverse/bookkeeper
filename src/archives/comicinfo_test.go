package archives

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseComicInfo(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="utf-8"?>
<ComicInfo xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
  <Title>The Amazing Spider-Man</Title>
  <Series>The Amazing Spider-Man</Series>
  <Number>1</Number>
  <Count>100</Count>
  <Volume>1</Volume>
  <AlternateSeries></AlternateSeries>
  <AlternateNumber>0</AlternateNumber>
  <StoryArc></StoryArc>
  <SeriesGroup></SeriesGroup>
  <AlternateCount>0</AlternateCount>
  <Summary>Peter Parker gets bitten by a radioactive spider.</Summary>
  <Notes>Scanned by ComicRack</Notes>
  <Year>1963</Year>
  <Month>3</Month>
  <Day>1</Day>
  <Writer>Stan Lee</Writer>
  <Penciller>Steve Ditko</Penciller>
  <Inker>Steve Ditko</Inker>
  <Colorist>Unknown</Colorist>
  <Letterer>Artie Simek</Letterer>
  <CoverArtist>Steve Ditko</CoverArtist>
  <Editor>Stan Lee</Editor>
  <Translator></Translator>
  <Publisher>Marvel Comics</Publisher>
  <Imprint></Imprint>
  <Genre>Superhero, Action</Genre>
  <Tags>classic, origin story</Tags>
  <Web></Web>
  <PageCount>20</PageCount>
  <LanguageISO>en</LanguageISO>
  <Format>Print</Format>
  <BlackAndWhite>No</BlackAndWhite>
  <Manga>No</Manga>
  <Characters>Spider-Man, J. Jonah Jameson</Characters>
  <Teams></Teams>
  <Locations>New York City</Locations>
  <ScanInformation></ScanInformation>
  <AgeRating>All Ages</AgeRating>
</ComicInfo>`

	bookInfo, err := parseComicInfo([]byte(xmlData))
	require.NoError(t, err, "should successfully parse ComicInfo.xml")

	// Test basic fields
	assert.Equal(t, "The Amazing Spider-Man", bookInfo.Title)
	assert.Equal(t, "The Amazing Spider-Man", bookInfo.Series)
	assert.Equal(t, "1", bookInfo.SeriesIndex)
	assert.Equal(t, "Peter Parker gets bitten by a radioactive spider.", bookInfo.Description)
	assert.Equal(t, "Marvel Comics", bookInfo.Publisher)
	assert.Equal(t, "1963-03-01", bookInfo.PublishedDate)
	assert.Equal(t, 20, bookInfo.Pages)
	assert.Equal(t, []string{"en"}, bookInfo.Language)

	// Test authors (should include all creators without duplicates)
	expectedAuthors := []string{"Stan Lee", "Steve Ditko", "Unknown", "Artie Simek"}
	assert.ElementsMatch(t, expectedAuthors, bookInfo.Authors)

	// Test keywords (should include genre, tags, characters, locations)
	expectedKeywords := []string{"Superhero", "Action", "classic", "origin story", "Spider-Man", "J. Jonah Jameson", "New York City"}
	assert.ElementsMatch(t, expectedKeywords, bookInfo.Keywords)
}

func TestParseComicInfoMinimal(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="utf-8"?>
<ComicInfo>
  <Series>Test Comic</Series>
  <Number>5</Number>
</ComicInfo>`

	bookInfo, err := parseComicInfo([]byte(xmlData))
	require.NoError(t, err, "should successfully parse minimal ComicInfo.xml")

	// Should construct title from series and number
	assert.Equal(t, "Test Comic #5", bookInfo.Title)
	assert.Equal(t, "Test Comic", bookInfo.Series)
	assert.Equal(t, "5", bookInfo.SeriesIndex)
}

func TestParseComicInfoDateFormats(t *testing.T) {
	tests := []struct {
		name     string
		year     int
		month    int
		day      int
		expected string
	}{
		{"Year only", 2023, 0, 0, "2023"},
		{"Year and month", 2023, 12, 0, "2023-12"},
		{"Full date", 2023, 12, 25, "2023-12-25"},
		{"Single digit month", 2023, 3, 0, "2023-03"},
		{"Single digit day", 2023, 12, 5, "2023-12-05"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xmlData := `<?xml version="1.0" encoding="utf-8"?>
<ComicInfo>
  <Title>Test</Title>`
			if tt.year > 0 {
				xmlData += fmt.Sprintf(`<Year>%d</Year>`, tt.year)
			}
			if tt.month > 0 {
				xmlData += fmt.Sprintf(`<Month>%d</Month>`, tt.month)
			}
			if tt.day > 0 {
				xmlData += fmt.Sprintf(`<Day>%d</Day>`, tt.day)
			}
			xmlData += `</ComicInfo>`

			bookInfo, err := parseComicInfo([]byte(xmlData))
			require.NoError(t, err)
			assert.Equal(t, tt.expected, bookInfo.PublishedDate)
		})
	}
}

func TestSplitCommaDelimited(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"Empty string", "", nil},
		{"Single item", "test", []string{"test"}},
		{"Multiple items", "a, b, c", []string{"a", "b", "c"}},
		{"With extra spaces", " a , b , c ", []string{"a", "b", "c"}},
		{"With empty parts", "a,, b, , c", []string{"a", "b", "c"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitCommaDelimited(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{"Empty slice", []string{}, []string{}},
		{"No duplicates", []string{"a", "b", "c"}, []string{"a", "b", "c"}},
		{"With duplicates", []string{"a", "b", "a", "c", "b"}, []string{"a", "b", "c"}},
		{"All same", []string{"a", "a", "a"}, []string{"a"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeDuplicates(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
