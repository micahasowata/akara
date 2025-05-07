package main

import (
	"bytes"
	"errors"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
	"github.com/yuin/goldmark"
)

func fileExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		return true, nil
	}

	return false, err
}

// converts the md file to html
func convertMdtoHTML(src string) (string, error) {
	exists, err := fileExist(src)
	if err != nil {
		return "", err
	}

	if !exists {
		return "", errors.New("source does not exist")
	}

	content, err := os.ReadFile(src)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = goldmark.Convert(content, &buf)
	if err != nil {
		return "", err
	}

	content = buf.Bytes()

	return string(content), nil
}

func createTargetPath(src string, target string) string {
	base := strings.TrimSuffix(filepath.Base(src), filepath.Ext(filepath.Base(src)))

	path := strings.Split(src, "/")
	parent := strings.Join(path[len(path)-2:len(path)-1], "")

	base = parent + "-" + base + ".html"

	return filepath.Join(target, base)
}

// // writes the content to the target
func createTargetDir(target string) error {
	info, err := os.Stat(target)
	if err == nil {
		if !info.IsDir() {
			return errors.New("target should be directory")
		}
	}

	if errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(target, 0o755)
		if err != nil {
			return err
		}
	}

	return err
}

// buildFile converts a markdown file into a corresponding html file
func build(layout, src, target string) error {
	content, err := convertMdtoHTML(src)
	if err != nil {
		return err
	}

	content, err = addLayout(layout, content)
	if err != nil {
		return err
	}

	err = createTargetDir(target)
	if err != nil {
		return err
	}

	targetPath := createTargetPath(src, target)

	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	content, err = m.String("text/html", content)
	if err != nil {
		return err
	}

	err = os.WriteFile(targetPath, []byte(content), 0o644)
	if err != nil {
		return err
	}

	return nil
}

func addLayout(layout string, content string) (string, error) {
	tmpl, err := template.ParseFiles(layout)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer

	err = tmpl.Execute(&buf, struct{ Content any }{Content: template.HTML(content)})
	if err != nil {
		return "", err
	}

	result := string(buf.Bytes())

	return result, nil
}
