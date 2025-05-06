package main

import (
	"log"
	"os"
)

func main() {
	err := buildFile("/home/micah/projects/akaratest/content/posts/post-10.md", "/home/micah/projects/akaratest/target")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
