package content

import (
	"bytes"
	"html/template"
)

type Layout struct {
	Content template.HTML
}

func NewLayout(c string) Layout {
	return Layout{
		Content: template.HTML(c),
	}
}

func ResolveContentLayout(src string, content Layout) ([]byte, error) {
	tmpl, err := template.ParseFiles(src)
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
