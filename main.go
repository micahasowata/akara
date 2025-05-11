package main

import (
	"fmt"

	"micahasowata.com/akara/content"
)

func main() {
	h, err := content.ConvertToHTML("/home/micah/projects/akaratest/content/posts/post-1.md")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(h))
}
