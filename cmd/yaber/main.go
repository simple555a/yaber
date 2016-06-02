package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/lmas/yaber"
)

var (
	prefix = flag.String("prefix", "asset_", "file prefix for the output files")
	strip  = flag.String("strip", "", "file path prefix to strip away from the asset files")
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("yaber: ")

	flag.Usage = usage
	flag.Parse()

	gen, e := yaber.NewGenerator(flag.Arg(1), "", *prefix, *strip)
	checkError(e)

	files, e := gen.GenerateAssets()
	checkError(e)

	for _, f := range files {
		e = ioutil.WriteFile(f.Path, f.Body, 0666)
		checkError(e)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `yaber v%s
Yet another binary embedder - Generate go code embedded with binary (and
gzip'ed) data of your local assets.

Usage:
  yaber [flags] /path/to/assets/dir/

Flags:
`, yaber.VERSION)
	flag.PrintDefaults()
}

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
