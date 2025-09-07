package archives

import (
	"fmt"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/webassembly"
)

var (
	pool     pdfium.Pool
	instance pdfium.Pdfium
)

func init() {
	// Initialize PDFium pool with WebAssembly
	var err error
	pool, err = webassembly.Init(webassembly.Config{
		MinIdle:  1, // Ensures at least 1 worker is always available
		MaxIdle:  1, // Limits to 1 worker when idle
		MaxTotal: 1, // Maximum 1 worker total
	})
	if err != nil {
		panic(fmt.Sprintf("failed to initialize PDFium pool: %v", err))
	}

	instance, err = pool.GetInstance(time.Second * 30)
	if err != nil {
		panic(fmt.Sprintf("failed to get PDFium instance: %v", err))
	}
}

func getBookInfoPDF(path string) (BookInfo, error) {
	// Load the PDF file into a byte array
	pdfBytes, err := os.ReadFile(path)
	if err != nil {
		return BookInfo{}, fmt.Errorf("failed to read PDF file: %w", err)
	}

	// Open the PDF using PDFium
	doc, err := instance.OpenDocument(&requests.OpenDocument{
		File: &pdfBytes,
	})
	if err != nil {
		return BookInfo{}, fmt.Errorf("failed to open PDF document: %w", err)
	}

	// Always close the document to release resources
	defer instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
		Document: doc.Document,
	})

	// Get page count
	pageCount, err := instance.FPDF_GetPageCount(&requests.FPDF_GetPageCount{
		Document: doc.Document,
	})
	if err != nil {
		return BookInfo{}, fmt.Errorf("failed to get page count: %w", err)
	}

	title := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	var authors []string
	var keywords []string

	// Get metadata
	metadata, err := instance.GetMetaData(&requests.GetMetaData{
		Document: doc.Document,
	})
	if err == nil && metadata != nil {
		for _, tag := range metadata.Tags {
			switch tag.Tag {
			case "Author":
				if tag.Value != "" {
					authors = []string{tag.Value}
				}
			case "Title":
				if tag.Value != "" && tag.Value != ".pdf" {
					title = tag.Value
				}
			case "Keywords":
				if tag.Value != "" {
					for _, s := range strings.Split(tag.Value, ",") {
						s = strings.TrimSpace(s)
						if s != "" {
							keywords = append(keywords, s)
						}
					}
				}
			}
		}
	}

	return BookInfo{
		Title:    title,
		Pages:    pageCount.PageCount,
		Authors:  authors,
		Keywords: keywords,
	}, nil
}

// extractPDF renders PDF pages as JPEG images using go-pdfium
func extractPDF(inputFile, outputFolder string) ([]Page, error) {
	// Load the PDF file into a byte array
	pdfBytes, err := os.ReadFile(inputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF file: %w", err)
	}

	// Open the PDF using PDFium
	doc, err := instance.OpenDocument(&requests.OpenDocument{
		File: &pdfBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open PDF document: %w", err)
	}

	// Always close the document to release resources
	defer instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
		Document: doc.Document,
	})

	// Get page count
	pageCount, err := instance.FPDF_GetPageCount(&requests.FPDF_GetPageCount{
		Document: doc.Document,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get page count: %w", err)
	}

	var pages []Page

	// Extract each page as JPEG
	for pageNum := 0; pageNum < pageCount.PageCount; pageNum++ {
		// Render page to image using go-pdfium
		renderPage, err := instance.RenderPageInDPI(&requests.RenderPageInDPI{
			Page: requests.Page{
				ByIndex: &requests.PageByIndex{
					Document: doc.Document,
					Index:    pageNum,
				},
			},
			DPI: 150, // Set DPI for good quality
		})
		if err != nil {
			return nil, fmt.Errorf("failed to render page %d: %w", pageNum+1, err)
		}

		// Get image dimensions from the rendered image
		width := renderPage.Result.Image.Bounds().Dx()
		height := renderPage.Result.Image.Bounds().Dy()

		// Generate filename with zero-padded page number
		filename := fmt.Sprintf("page_%02d.jpg", pageNum+1)
		outputPath := filepath.Join(outputFolder, filename)

		// Create output file
		file, err := os.Create(outputPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create output file %s: %w", outputPath, err)
		}

		// Encode as JPEG
		err = jpeg.Encode(file, renderPage.Result.Image, &jpeg.Options{Quality: 90})
		file.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to encode page %d as JPEG: %w", pageNum+1, err)
		}

		// Add to pages list with dimensions
		pages = append(pages, Page{
			Path:   filename,
			Width:  width,
			Height: height,
		})
	}

	return pages, nil
}
