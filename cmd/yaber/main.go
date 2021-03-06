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
	pkg    = flag.String("pkg", "", "package name to use for the generated code")
	output = flag.String("out", "assets", "file path prefix for the generated files")
	strip  = flag.String("strip", "", "file path prefix to strip away from the assets")
	public = flag.Bool("public", false, "export public functions for getting assets")
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("yaber: ")

	flag.Usage = usage
	flag.Parse()

	gen, e := yaber.NewGenerator(*pkg, *output, *strip, *public)
	checkError(e)

	files, e := gen.Generate(flag.Args())
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
  yaber [flags] /paths/to/assets/dirs/

Flags:
`, yaber.VERSION)
	flag.PrintDefaults()
}

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
