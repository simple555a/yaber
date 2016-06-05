package main

import (
	"bytes"
	"fmt"

	"github.com/lmas/yaber/example/assets"
)

func main() {
	fmt.Println("Comparing 'templates/hello' file...")
	embed, e := assets.Asset("templates/hello")
	checkError(e)
	assets.SetRawAssets(true)
	raw, e := assets.Asset("templates/hello")
	checkError(e)
	if bytes.Equal(embed, raw) == true {
		fmt.Println("Both versions match.")
	} else {
		fmt.Printf("FILES NOT MATCHED! embed='%s', raw='%s'\n", embed, raw)
	}

	fmt.Println("Comparing all files in 'templates/'...")
	rawfiles, e := assets.AssetDir("templates/")
	checkError(e)
	assets.SetRawAssets(false)
	embedfiles, e := assets.AssetDir("templates/")
	checkError(e)
	if mapsEqual(embedfiles, rawfiles) == true {
		fmt.Println("Both versions match.")
	} else {
		fmt.Println("DIRS NOT MATCHED!")
	}
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func mapsEqual(m1, m2 map[string][]byte) bool {
	if len(m1) != len(m2) {
		return false
	}

	for k, v := range m1 {
		b, ok := m2[k]
		if !ok {
			return false
		}
		if !bytes.Equal(v, b) {
			return false
		}
	}
	return true
}
