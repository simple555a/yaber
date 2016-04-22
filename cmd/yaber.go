package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/lmas/yaber"
)

func main() {
	app := cli.NewApp()
	app.Version = yaber.VERSION
	app.HideHelp = true // hides the help command, eww
	app.Usage = "Generate go code with embedded binary data from asset files"
	app.ArgsUsage = "/path/to/assets/dir/"
	app.Flags = []cli.Flag{
		cli.BoolFlag{ // Have to add it here so it will show up in global options
			Name:  "help",
			Usage: "show help",
		},
		cli.StringFlag{
			Name:  "pkg",
			Usage: "package name to use in generated files",
			Value: "assets",
		},
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
