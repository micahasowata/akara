package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
)

// buildFile converts a markdown file into a corresponding html file
func buildFile(src string, target string) error {
	content, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = goldmark.Convert(content, &buf)
	if err != nil {
		return err
	}

	info, err := os.Stat(target)
	if os.IsNotExist(err) {
		err := os.Mkdir(target, 0755)
		if err != nil {
			return err
		}
	}

	if !info.IsDir() {
		return errors.New("target is not a directory")
	}

	ext := filepath.Ext(filepath.Base(src))

	targetBase := strings.TrimSuffix(filepath.Base(src), ext) + ".html"

	targetPath := filepath.Join(target, targetBase)

	err = os.WriteFile(targetPath, buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}
