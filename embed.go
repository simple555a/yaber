package yaber

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"
)

// Make a dev template for loading files from disk and returns it as go code
// in []byte format.
func MakeDevAsset(pkgName string) ([]byte, error) {
	data := map[string]interface{}{
		"pkgName": pkgName,
		"ver":     VERSION,
	}
	tmpl := runTemplate(tmplDevAsset, data)
	return tmpl, nil
}

// Make a build template for loading embedded files, when running with the build
// tag "embed".
// Output is the same as MakeDevAsset: go code in []byte format.
func MakeBuildAsset(pkgName, path string, stripPath string) ([]byte, error) {
	files, e := embedAssets(path, stripPath)
	if e != nil {
		return nil, e
	}
	data := map[string]interface{}{
		"pkgName":  pkgName,
		"fileData": files,
		"ver":      VERSION,
	}
	tmpl := runTemplate(tmplBuildAsset, data)
	return tmpl, nil
}

// Compile the tmpl string, executes it using data and returns the result.
func runTemplate(tmpl string, data map[string]interface{}) []byte {
	t := template.Must(template.New("").Parse(tmpl))
	buf := new(bytes.Buffer)
	t.Execute(buf, data)
	return buf.Bytes()
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
