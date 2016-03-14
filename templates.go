package yaber

var tmplDevAsset = `// +build !embed

package {{.pkgName}}

import "io/ioutil"

func Asset(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}
`

var tmplBuildAsset = `// +build embed

package {{.pkgName}}

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"os"
)

func Asset(path string) ([]byte, error) {
	body, ok := _rawAssets[path]
	if !ok {
		return nil, os.ErrNotExist
	}

	buf := bytes.NewBuffer(body)
	gr, e := gzip.NewReader(buf)
	if e != nil {
		return nil, e
	}
	defer gr.Close()

	return ioutil.ReadAll(gr)
}

var _rawAssets = map[string][]byte{
{{range $path, $body := .fileData}}
	"{{$path}}": []byte({{printf "%+q" $body}}),
{{end -}}
}
`

// The "{{end -}}
// }
// `" part is a bit of a hack to ensure we get a correctly gofmt'ed output file
// and makes git shut up about "No newline at end of file".
