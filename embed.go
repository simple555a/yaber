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
	FilePath   string
	Package    string
	OutputFile string
	StripPath  string
}

func NewGenerator(path, pkg, output, strip string) (*Generator, error) {
	if len(path) < 1 {
		return nil, ErrNoPath
	}

	if len(output) < 1 {
		output = "assets.go"
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
		FilePath:   path,
		Package:    pkg,
		OutputFile: output,
		StripPath:  strip,
	}
	return g, nil
}

func (g *Generator) GenerateAssets() ([]*AssetFile, error) {
	files, e := embedAssets(g.FilePath, g.StripPath)
	if e != nil {
		return nil, e
	}

	data := map[string]interface{}{
		"version": VERSION,
		"package": g.Package,
		"command": executedCommand(),
		"files":   files,
	}

	// Generate the *dev.go file with no assets.
	mainBody, e := runTemplate(tmplMain, data)
	if e != nil {
		return nil, e
	}
	main := &AssetFile{
		Path: g.OutputFile,
		Body: mainBody,
	}

	return []*AssetFile{main}, nil
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
