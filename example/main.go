package main

import (
	"fmt"

	"github.com/lmas/yaber/example/assets"
)

func main() {
	files := []string{
		"templates/aaa.txt",
		"templates/empty_file",
		"templates/hello",
	}

	for _, p := range files {
		body, e := assets.Asset(p)
		if e != nil {
			panic(e)
		}
		fmt.Printf("file %s (%d bytes)\n", p, len(body))
		fmt.Println(string(body))
	}
}
