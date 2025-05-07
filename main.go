package main

import (
	"log"
	"os"
)

func main() {
	err := build("/home/micah/projects/akaratest/layouts/base.html", "/home/micah/projects/akaratest/content/posts/post-2.md", "/home/micah/projects/akaratest/target")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
