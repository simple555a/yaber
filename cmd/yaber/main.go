package main

import (
	"io/ioutil"
	"log"
	"os"

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
	path := c.Args().First()
	prefix := c.GlobalString("prefix")
	strip := c.GlobalString("strip")

	gen, e := yaber.NewGenerator(path, "", prefix, strip)
	checkError(e)

	files, e := gen.GenerateAssets()
	checkError(e)

	for _, f := range files {
		e = ioutil.WriteFile(f.Path, f.Body, 0666)
		checkError(e)
	}
}

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
