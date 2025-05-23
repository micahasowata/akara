package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/dlclark/regexp2"
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
	err = os.MkdirAll(target, 0o777)
	if err != nil {
		return err
	}

	wd := walkerDir{
		src:    src,
		target: target,
	}

	err = os.MkdirAll(target, 0o755)
	if err != nil {
		return err
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

	if entry.IsDir() {
		_, err = os.Stat(filepath.Join(path, "_index.html"))
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("content and its subdirectories must have _index.html page")
		}
	}

	bpath := strings.Split(path, "/")

	if filepath.Ext(entry.Name()) == ".html" && entry.Name() != "_index.html" && bpath[len(bpath)-2] != "content" {
		return errors.New("only content directory can have multiple standalone pages")
	}

	relPath, err := filepath.Rel(wd.src, path)
	if err != nil {
		return err
	}

	relPath = filepath.Clean(strings.TrimPrefix(relPath, "content"))

	targetPath := filepath.Join(wd.target, relPath)

	if entry.IsDir() {
		err = os.MkdirAll(targetPath, fs.ModePerm)
		if err != nil {
			return err
		}
	}

	if filepath.Ext(entry.Name()) == ".html" {
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		err = os.WriteFile(targetPath, content, fs.ModePerm)
		if err != nil {
			return err
		}
	}

	if filepath.Ext(entry.Name()) == ".md" {
		layoutPath := filepath.Join(wd.src, "layouts", strings.Split(relPath, "/")[1]+".html")

		tgPath := filepath.Join(wd.target, filepath.Clean(strings.TrimSuffix(relPath, filepath.Ext(relPath))+".html"))

		content, err := convertToHTML(path)
		if err != nil {
			return err
		}

		content, err = resolveLayout(layoutPath, layout{Content: template.HTML(string(content))})
		if err != nil {
			return err
		}

		err = os.WriteFile(tgPath, content, fs.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func fileExist(path string) bool {
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}

	if info.IsDir() {
		return false
	}

	return true
}

func dirExist(path string) bool {
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}

	if !info.IsDir() {
		return false
	}

	return true
}

func getLinksToStylesheets(content []byte) ([]string, error) {
	start := bytes.Index(content, []byte("<head>")) + len("<head>")
	end := bytes.Index(content, []byte("</head>"))

	focusArea := content[start:end]

	pattern, err := regexp2.Compile(`<link(?=[^>]*\srel\s*=\s*(['"])stylesheet\1)(?:[^>]*?\s)?href\s*=\s*(['"])((?!https:\/\/).*?)\2[^>]*?\/?>`, 0)
	if err != nil {
		return nil, nil
	}

	var matches []string

	m, err := pattern.FindStringMatch(string(focusArea))
	if err != nil {
		return nil, err
	}

	for m != nil {
		matches = append(matches, m.String())

		m, err = pattern.FindNextMatch(m)
		if err != nil {
			return nil, err
		}
	}

	result, err := pattern.Replace(string(focusArea), "", -1, -1)
	if err != nil {
		return nil, err
	}

	content = slices.Replace(content, start, end, []byte(result)...)

	fmt.Println(string(content))

	return matches, nil
}

func findStyleSheetPaths(links []string) ([]string, error) {
	var paths []string

	for _, v := range links {
	}
	return paths, nil
}

func serveFiles(w http.ResponseWriter, r *http.Request, target string) error {
	path := filepath.Clean(r.URL.Path)
	sParts := strings.Split(path, "/")

	// root path
	if len(sParts) == 2 {
		// homepage
		if sParts[1] == "" {
			pPath := filepath.Join(target, "_index.html")
			if fileExist(pPath) {
				http.ServeFile(w, r, pPath)
			}
		}

		if sParts[1] != "" {
			ok := dirExist(filepath.Join(target, sParts[1]))
			if ok {
				http.ServeFile(w, r, filepath.Join(target, sParts[1], "_index.html"))
			} else {
				if fileExist(filepath.Join(target, "_"+sParts[1]+".html")) {
					http.ServeFile(w, r, filepath.Join(target, "_"+sParts[1]+".html"))
				}
			}
		}
	}

	// other index paths
	ok := dirExist(filepath.Join(target, path))
	if ok {
		http.ServeFile(w, r, filepath.Join(target, path, "_index.html"))
	} else {
		if fileExist(filepath.Join(target, path+".html")) {
			http.ServeFile(w, r, filepath.Join(target, path+".html"))
		}
	}
	return nil
}

func main() {
	// err := buildContent("/home/micah/projects/akaratest")
	// if err != nil {
	// 	log.Println(err.Error())
	// }

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	err := serveFiles(w, r, "/home/micah/projects/akaratest/target")
	// 	if err != nil {
	// 		fmt.Fprint(w, err.Error())
	// 	}
	// })

	// http.ListenAndServe(":8080", nil)
	//
	htm, err := convertToHTML("/home/micah/projects/akaratest/content/posts/posts-1.md")
	if err != nil {
		panic(err)
	}

	htm, err = resolveLayout("/home/micah/projects/akaratest/layouts/posts.html", layout{Content: template.HTML(htm)})
	if err != nil {
		panic(err)
	}

	links, err := getLinksToStylesheets(htm)
	if err != nil {
		panic(err)
	}

	if links != nil {
		fmt.Println(links)
	}
}
