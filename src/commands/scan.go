package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/biblioteca/bookkeeper/src/archives"
	"github.com/biblioteca/bookkeeper/src/utils"
)

type metadata struct {
	Path   string            `json:"path"`
	Status string            `json:"status"`
	Size   int64             `json:"size,omitempty"`
	Hash   string            `json:"hash"`
	Book   archives.BookInfo `json:"book"`
}

type errorMetadata struct {
	Path   string `json:"path"`
	Status string `json:"status"`
	Error  string `json:"error"`
}

func Scan(scanPath string) error {
	abs, err := filepath.Abs(scanPath)
	if err != nil {
		return err
	}

	return filepath.WalkDir(abs, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			fmt.Println(errorLine(abs, path, err))
			return nil
		}
		if d.IsDir() {
			return nil
		}
		if !utils.IsValidBookFile(path) {
			return nil
		}

		line, err := scanBook(abs, path)
		if err != nil {
			fmt.Println(errorLine(abs, path, err))
			return nil
		}
		fmt.Println(line)
		return nil
	})
}

func scanBook(root string, path string) (string, error) {
	book, err := archives.GetBookInfo(path)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	m := metadata{
		Path:   relPath(root, path),
		Status: "success",
		Hash:   "",
		Size:   info.Size(),
		Book:   book,
	}
	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func errorLine(root, path string, e error) string {
	m := errorMetadata{
		Path:   relPath(root, path),
		Status: "failed",
		Error:  e.Error(),
	}
	b, _ := json.Marshal(m)
	return string(b)
}

func relPath(root, p string) string {
	rel, err := filepath.Rel(root, p)
	if err != nil {
		return p
	}
	return rel
}
