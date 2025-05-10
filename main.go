package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"micahasowata.com/akara/build"
)

func main() {
	html, err := build.ConvertMdToHTML("/home/micah/projects/akaratest/content/posts/post-10.md")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	html, err = build.ResolveToLayout("/home/micah/projects/akaratest/layouts/base.html", html)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	html, err = build.RemoveLocalStyleLinks(html)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	localStyles := build.ReadInFileStyles(html)
	extStyles, err := build.ReadSectionStyles("/home/micah/projects/akaratest/css/posts.css")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	ogCSS := append(extStyles, localStyles...)

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	fmt.Println(string(ogCSS))

	mCSS, err := m.Bytes("text/css", ogCSS)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Println(string(mCSS))
}
