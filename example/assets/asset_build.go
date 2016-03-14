// +build embed

package assets

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

	"templates/aaa.txt": []byte("\x1f\x8b\b\x00\x00\tn\x88\x00\xff\n\xc9\xc8,VH\xcb\xccIU\x00\xd2Y\xa5\xc5%\n\x89\n%\xa9\xc5%z\\\x00\x00\x00\x00\xff\xff\x01\x00\x00\xff\xff\x9d\xc5\x12$\x1a\x00\x00\x00"),

	"templates/empty_file": []byte("\x1f\x8b\b\x00\x00\tn\x88\x00\xff\x00\x00\x00\xff\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00"),

	"templates/hello": []byte("\x1f\x8b\b\x00\x00\tn\x88\x00\xff\xf2H\xcd\xc9\xc9W(\xcf/\xcaIQ\xe4\x02\x00\x00\x00\xff\xff\x01\x00\x00\xff\xffA\u4a72\r\x00\x00\x00"),
}
