package yaber

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var ErrNoPath = errors.New("no file path to assets")

type AssetFile struct {
	Path string
	Body []byte
}

type Generator struct {
	FilePaths   []string
	Package     string
	OutputFile  string
	StripPath   string
	PublicFuncs bool
}

func NewGenerator(path []string, pkg, output, strip string, publicFuncs bool) (*Generator, error) {
	if len(path) < 1 {
		return nil, ErrNoPath
	}

	if len(output) < 1 {
		output = "assets"
	}

	if len(pkg) < 1 {
		var e error
		// Default to use the output (or the current) dir as the pkg name
		pkg, e = getPackageName(filepath.Dir(output))
		if e != nil {
			return nil, e
		}
	}

	// TODO: support multiple file paths/dirs
	g := &Generator{
		FilePaths:   path,
		Package:     pkg,
		OutputFile:  output,
		StripPath:   strip,
		PublicFuncs: publicFuncs,
	}
	return g, nil
}

func (g *Generator) GenerateAssets() ([]*AssetFile, error) {
	files := make(map[string][]byte)
	for _, p := range g.FilePaths {
		f, e := embedAssets(p, g.StripPath)
		if e != nil {
			return nil, e
		}
		for k, v := range f {
			files[k] = v
		}
	}

	data := map[string]interface{}{
		"version": VERSION,
		"package": g.Package,
		"command": executedCommand(),
		"files":   files,
	}
	if g.PublicFuncs {
		data["assetFunc"] = "Asset"
		data["setRawFunc"] = "SetRawAssets"
	} else {
		data["assetFunc"] = "asset"
		data["setRawFunc"] = "setRawAssets"
	}

	// Generate the main file with embedded files.
	mainBody, e := runTemplate(tmplMain, data)
	if e != nil {
		return nil, e
	}
	main := &AssetFile{
		Path: g.OutputFile + ".go",
		Body: mainBody,
	}

	// Generate the test file.
	var first string
	for k, _ := range files {
		first = k
		break
	}
	data["firstPath"] = first
	data["firstBody"] = files[first]
	data["dirs"] = g.FilePaths

	testBody, e := runTemplate(tmplTest, data)
	if e != nil {
		return nil, e
	}
	test := &AssetFile{
		Path: g.OutputFile + "_test.go",
		Body: testBody,
	}

	return []*AssetFile{main, test}, nil
}

// Recursively reads all regular files in path, into memory as gzipped data.
// Returns a map where the keys are file paths and the values are the gzip byte data.
func embedAssets(path string, stripPath string) (map[string][]byte, error) {
	list := make(map[string][]byte)
	dirs := []string{path}

	for len(dirs) > 0 {
		d := dirs[0]
		dirs = dirs[1:]
		files, e := ioutil.ReadDir(d)
		if e != nil {
			return nil, e
		}

		for _, f := range files {
			fpath := filepath.Join(d, f.Name())
			tmpPath := strings.TrimPrefix(fpath, stripPath)

			if f.IsDir() {
				dirs = append(dirs, fpath)
				continue
			}
			if !f.Mode().IsRegular() {
				continue
			}

			fbody, e := ioutil.ReadFile(fpath)
			if e != nil {
				return nil, e
			}
			if len(fbody) < 1 {
				list[tmpPath] = []byte{}
				continue
			}

			buf := new(bytes.Buffer)
			gw := gzip.NewWriter(buf)
			defer gw.Close()

			if _, e = gw.Write(fbody); e != nil {
				return nil, e
			}
			gw.Flush()
			if gw.Close() != nil {
				return nil, e
			}

			list[tmpPath] = buf.Bytes()
		}
	}
	return list, nil
}
