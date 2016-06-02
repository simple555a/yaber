package main

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/lmas/yaber"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("yaber: ")

	app := cli.NewApp()
	app.Version = yaber.VERSION
	app.Usage = "Generate go code with embedded binary data from asset files"
	app.ArgsUsage = "/path/to/assets/dir/"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "prefix",
			Usage: "file prefix for generated files",
			Value: "asset_",
		},
		cli.StringFlag{
			Name:  "strip",
			Usage: "file path prefix to strip away",
		},
	}
	app.Action = generateFiles
	app.Run(os.Args)
}

func generateFiles(c *cli.Context) {
	path := strings.TrimSpace(c.Args().First())
	if len(path) < 1 {
		fmt.Println("Error: No file path specified.")
		os.Exit(1)
	}
	prefix := c.GlobalString("prefix")
	strip := c.GlobalString("strip")

	// Default to use the prefix (or the current) dir as the pkg name
	dir, e := filepath.Abs(filepath.Dir(prefix))
	checkError(e)
	pkgName := filepath.Base(dir)

	// Check if the dir contains a go pkg
	pkg, e := build.ImportDir(dir, 0)
	if e != nil {
		// NoGoError = not a valid pkg in dir, so ignore it and use the
		// default from above.
		if _, ok := e.(*build.NoGoError); !ok {
			checkError(e)
		}
	} else {
		// Looks like it was an actual go pkg, so use that
		pkgName = pkg.Name
	}

	dev, e := yaber.MakeDevAsset(pkgName)
	checkError(e)
	build, e := yaber.MakeBuildAsset(pkgName, path, strip)
	checkError(e)

	dp := fmt.Sprintf("%sdev.go", prefix)
	writeAsset(dp, dev)
	bp := fmt.Sprintf("%sbuild.go", prefix)
	writeAsset(bp, build)
}

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func joinPath(path string) string {
	cp, e := os.Getwd()
	checkError(e)
	return filepath.Join(cp, path)
}

func writeAsset(path string, data []byte) {
	p := joinPath(path)
	e := ioutil.WriteFile(p, data, 0666)
	checkError(e)
}
