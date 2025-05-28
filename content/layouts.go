package content

import (
	"bytes"
	"html/template"
)

// Layout holds the HTML content that would be injected into a page
type Layout struct {
	Content template.HTML
}

// NewLayout creates a new populated instance of Layout
func NewLayout(content []byte) Layout {
	return Layout{
		Content: template.HTML(string(content)),
	}
}

// ParseContent parses the content into the layout
func ParseContent(layout string, content Layout) ([]byte, error) {
	tmpl, err := template.ParseFiles(layout)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, content)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
