package main

import (
	"fmt"

	"github.com/lmas/yaber/example/assets"
)

func main() {
	fmt.Println("Get a specific file.")
	path := "templates/hello"
	body, e := assets.Asset(path)
	if e != nil {
		panic(e)
	}
	fmt.Printf("file %s (%d bytes)\n", path, len(body))
	fmt.Println(string(body))

	fmt.Println("Get all files with a specific path prefix (or no prefix).")
	files, e := assets.AssetDir("templates/")
	if e != nil {
		panic(e)
	}

	for path, body := range files {
		fmt.Printf("file %s (%d bytes)\n", path, len(body))
		fmt.Println(string(body))
	}
}
