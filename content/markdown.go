package content

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

// isMd checks if a file path ends in .md
func isMd(src string) bool {
	ext := filepath.Ext(src)

	return ext == ".md"
}

// ConvertToHTML transforms the content of a file into HTML
//
// src should be to a .md file
func ConvertToHTML(src string) ([]byte, error) {
	if !isMd(src) {
		return nil, errors.New("src must be path to a markdown file")
	}

	md := goldmark.New(goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithHeadingAttribute(),
		),
	)

	content, err := os.ReadFile(src)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	err = md.Convert(content, &buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
