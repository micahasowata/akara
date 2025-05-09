package build

import (
	"bufio"
	"bytes"
	"fmt"
	"html"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

// ConvertMdToHTML converts md to HTML
//
//	it takes the path
func ConvertMdToHTML(src string) ([]byte, error) {
	file, err := os.Open(src)
	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(file.Name())
	if ext != ".md" {
		return nil, Error{fmt.Sprintf("source file should be markdown. it is %s", strings.ToLower(ext))}
	}

	reader := bufio.NewReader(file)

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	converter := goldmark.New(goldmark.WithExtensions(extension.GFM))

	var buf bytes.Buffer

	err = converter.Convert(content, &buf)
	if err != nil {
		return nil, err
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type Layout struct {
	Content template.HTML
}

func ResolveToLayout(layout string, content []byte) ([]byte, error) {
	tmpl, err := template.ParseFiles(layout)
	if err != nil {
		return nil, nil
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, Layout{Content: template.HTML(html.UnescapeString(string(content)))})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
