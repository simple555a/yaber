package yaber

import "text/template"

var (
	tmplDev   *template.Template
	tmplBuild *template.Template
)

func init() {
	head := template.Must(template.New("").Parse(rawHead))

	tmplDev = template.Must(head.Clone())
	template.Must(tmplDev.Parse(rawDev))

	tmplBuild = template.Must(head.Clone())
	template.Must(tmplBuild.Parse(rawBuild))
}

var rawHead = `{{define "head"}}
//go:generate {{.command}}

// Code generated by yaber v{{.version}} (https://github.com/lmas/yaber)
// DO NOT EDIT.

package {{.package}}{{end}}`

var rawDev = `// +build !{{.tag}}
{{template "head" .}}

import (
        "io/ioutil"
        "path/filepath"
)

func Asset(path string) ([]byte, error) {
        return ioutil.ReadFile(path)
}

func AssetDir(dir string) (map[string][]byte, error) {
        list := make(map[string][]byte)
        dirs := []string{dir}

        for len(dirs) > 0 {
                d := dirs[0]
                dirs = dirs[1:]
                files, e := ioutil.ReadDir(d)
                if e != nil {
                        return nil, e
                }

                for _, f := range files {
                        fpath := filepath.Join(d, f.Name())

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
                        list[fpath] = fbody
                }
        }
        return list, nil
}
`

var rawBuild = `// +build {{.tag}}
{{template "head" .}}

import (
        "bytes"
        "compress/gzip"
        "io/ioutil"
        "os"
        "strings"
)

func Asset(path string) ([]byte, error) {
        body, ok := _rawAssets[path]
        if !ok {
                return nil, os.ErrNotExist
        }
        return decompress(body)
}

func AssetDir(dir string) (map[string][]byte, error) {
        var e error
        files := make(map[string][]byte)
        for path, body := range _rawAssets {
                if strings.HasPrefix(path, dir) {
                        files[path], e = decompress(body)
                        if e != nil {
                                return nil, e
                        }
                }
        }
        return files, nil
}

func decompress(data []byte) ([]byte, error) {
        buf := bytes.NewBuffer(data)
        gr, e := gzip.NewReader(buf)
        if e != nil {
                return nil, e
        }
        defer gr.Close()
        return ioutil.ReadAll(gr)
}

var _rawAssets = map[string][]byte{
{{range $path, $body := .files}}
	"{{$path}}": []byte({{printf "%+q" $body}}),
{{end -}}
}
`

// The "{{end -}}
// }
// `" part is a bit of a hack to ensure we get a correctly gofmt'ed output file
// and makes git shut up about "No newline at end of file".
