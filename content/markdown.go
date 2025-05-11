package content

import (
	"bytes"
	"os"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

func ConvertToHTML(src string) ([]byte, error) {
	content, err := os.ReadFile(src)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	md := goldmark.New(goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	err = md.Convert(content, &buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
