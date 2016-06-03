package yaber

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
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
	FilePath     string
	Package      string
	OutputPrefix string
	StripPath    string
	BuildTag     string
}

func NewGenerator(path, pkg, prefix, strip, tag string) (*Generator, error) {
	if len(path) < 1 {
		return nil, ErrNoPath
	}

	if len(prefix) < 1 {
		prefix = "asset_"
	}

	if len(pkg) < 1 {
		var e error
		// Default to use the prefix (or the current) dir as the pkg name
		pkg, e = getPackageName(filepath.Dir(prefix))
		if e != nil {
			return nil, e
		}
	}

	if len(tag) < 1 {
		tag = "embed"
	}

	// TODO: support multiple file paths/dirs
	g := &Generator{
		FilePath:     path,
		Package:      pkg,
		OutputPrefix: prefix,
		StripPath:    strip,
		BuildTag:     tag,
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
		"tag":     g.BuildTag,
		"files":   files,
	}

	// Generate the *dev.go file with no assets.
	devBody, e := runTemplate(tmplDev, data)
	if e != nil {
		return nil, e
	}
	dev := &AssetFile{
		Path: fmt.Sprintf("%sdev.go", g.OutputPrefix),
		Body: devBody,
	}

	// Generate the *build.go file with embedded assets.
	buildBody, e := runTemplate(tmplBuild, data)
	if e != nil {
		return nil, e
	}

	build := &AssetFile{
		Path: fmt.Sprintf("%sbuild.go", g.OutputPrefix),
		Body: buildBody,
	}

	return []*AssetFile{dev, build}, nil
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
