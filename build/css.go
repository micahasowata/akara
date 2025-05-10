package build

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"slices"

	"github.com/dlclark/regexp2"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
)

func RemoveLocalStyleLinks(html []byte) ([]byte, error) {
	start := bytes.Index(html, []byte("<head>")) + len("<head>")
	end := bytes.Index(html, []byte("</head>"))

	focusArea := html[start:end]

	pattern, err := regexp2.Compile(`<link(?=[^>]*\srel\s*=\s*(['"])stylesheet\1)(?:[^>]*?\s)?href\s*=\s*(['"])((?!https:\/\/).*?)\2[^>]*?\/?>`, 0)
	if err != nil {
		return nil, err
	}

	result, err := pattern.Replace(string(focusArea), "", -1, -1)
	if err != nil {
		return nil, err
	}

	return slices.Replace(html, start, end, []byte(result)...), nil
}

func ReadSectionStyles(src string) ([]byte, error) {
	file, err := os.Open(src)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(file)
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func ReadInFileStyles(html []byte) []byte {
	start := bytes.Index(html, []byte("<style>")) + len("<style>")
	end := bytes.Index(html, []byte("</style>"))

	return html[start:end]
}

func MergeStyleSheets(local, external []byte) []byte {
	m := minify.New()
	m.Add("text/css", css.Minify(m *minify.M, w io.Writer, r io.Reader, params map[string]string))
}
