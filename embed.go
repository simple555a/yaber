package yaber

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// ErrNoPaths is returned from NewGenerator() whenever the user has failed to
// provide at least a single file path to assets.
var ErrNoPaths = errors.New("no file paths to assets")

// AssetFile is the final, generated product from a Generator.
type AssetFile struct {
	Path string
	Body []byte
}

// Generator is the main object used for generating new files with embedded
// assets and tests.
type Generator struct {
	// FilePaths is the list of asset directories/files to embed.
	FilePaths []string

	// Package sets the package name for the newly generated files.
	Package string

	// OutputFile is the path prefix to append to the generated files.
	OutputFile string

	// StripPath will strip this prefix from the embedded asset file paths.
	StripPath string

	// If true, PublicFuncs let's you publicly export some functions for
	// accessing your embedded assets, from outside the package.
	PublicFuncs bool
}

// NewGenerator constructs a new Generator object.
func NewGenerator(path []string, pkg, output, strip string, publicFuncs bool) (*Generator, error) {
	if len(path) < 1 {
		return nil, ErrNoPaths
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

// GenerateAssets attempts to read the provided asset files, compress them and
// then embedd them in a new Go file, along with a basic test file.
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
	for k := range files {
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
