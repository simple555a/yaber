package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/lmas/yaber"
)

func main() {
	app := cli.NewApp()
	app.Version = yaber.VERSION
	app.Usage = "Generate go code with embedded binary data from asset files"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "path",
			Usage: "Path to asset directory for generating files from (REQUIRED)",
		},
		cli.StringFlag{
			Name:  "pkg",
			Usage: "Package name to use in generated files",
			Value: "assets",
		},
		cli.StringFlag{
			Name:  "prefix",
			Usage: "File prefix for generated files",
			Value: "asset_",
		},
		cli.StringFlag{
			Name:  "strip",
			Usage: "File path prefix to strip away",
		},
	}
	app.Action = generateFiles
	app.Run(os.Args)
}

func generateFiles(c *cli.Context) {
	path := c.GlobalString("path")
	pkgName := c.GlobalString("pkg")
	prefix := c.GlobalString("prefix")
	strip := c.GlobalString("strip")

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
		fmt.Println(e)
		os.Exit(1)
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
