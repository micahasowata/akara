package main

import (
	"bytes"
	"errors"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

type layout struct {
	Content template.HTML
}

func resolveLayout(src string, content layout) ([]byte, error) {
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

func convertToHTML(src string) ([]byte, error) {
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

func buildContent(src string) error {
	// determine project root by checking if it has the content and layouts directories
	info, err := os.Stat(filepath.Join(src, "content"))
	if errors.Is(err, os.ErrNotExist) {
		return err
	}

	if !info.IsDir() {
		return errors.New("content must be a directory")
	}

	info, err = os.Stat(filepath.Join(src, "layouts"))
	if errors.Is(err, os.ErrNotExist) {
		return err
	}

	if !info.IsDir() {
		return errors.New("layouts must be a directory")
	}

	target := filepath.Join(src, "target")
	err = os.MkdirAll(target, 0777)
	if err != nil {
		return err
	}

	wd := walkerDir{
		src:    src,
		target: target,
	}

	err = filepath.WalkDir(filepath.Join(src, "content"), wd.walker)
	if err != nil {
		return err
	}

	return nil
}

type walkerDir struct {
	src    string
	target string
}

func (wd walkerDir) walker(path string, entry fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	relPath, err := filepath.Rel(wd.src, path)
	if err != nil {
		return err
	}

	if !entry.IsDir() {
		relPath1 := filepath.Clean(strings.TrimPrefix(relPath, "content"))
		relPath2 := filepath.Clean(strings.Trim(relPath1, filepath.Ext(relPath1)))
		relPath3 := filepath.Clean(relPath2 + ".html")
		targetPath := filepath.Join(wd.target, relPath3)

		parts := strings.Split(relPath1, "/")
		layoutPath := filepath.Join(wd.src, "layouts", parts[1]+".html")

		info, err := os.Stat(layoutPath)
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("no layout file for content section")
		}

		if info.IsDir() {
			return errors.New("content layouts must be files")
		}

		err = os.MkdirAll(filepath.Dir(targetPath), 0755)
		if err != nil {
			return err
		}

		content, err := convertToHTML(path)
		if err != nil {
			return err
		}

		content, err = resolveLayout(layoutPath, layout{Content: template.HTML(string(content))})
		if err != nil {
			return err
		}

		err = os.WriteFile(targetPath, content, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func serveFiles(w http.ResponseWriter, r *http.Request, target string) error {
	cp := filepath.Clean(r.URL.Path)
	tp := filepath.Join(target, cp+".html")

	info, err := os.Stat(tp)
	if errors.Is(err, os.ErrNotExist) {
		return err
	}

	if info.IsDir() {
		return errors.New("file should be a file")
	}

	http.ServeFile(w, r, tp)
	return nil
}

func main() {
	// err := buildContent("/home/micah/projects/akaratest")
	// if err != nil {
	// 	panic(err)
	// }

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := serveFiles(w, r, "/home/micah/projects/akaratest/target")
		if err != nil {
			log.Fatal(err)
		}
	})

	http.ListenAndServe(":8080", nil)
}
