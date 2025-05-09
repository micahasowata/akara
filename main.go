package main

import (
	"fmt"
	"log"
	"os"

	"micahasowata.com/akara/build"
)

func main() {
	html, err := build.ConvertMdToHTML("/home/micah/projects/akaratest/content/posts/post-1.md")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	html, err = build.ResolveToLayout("/home/micah/projects/akaratest/layouts/base.html", html)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Println(string(html))
}
