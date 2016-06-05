package yaber

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"testing"
)

func assert(t *testing.T, val, expected interface{}) {
	if val != expected {
		t.Errorf("Expected: %#v, Got: %#v\n", expected, val)
	}
}

func TestEmbedAssets(t *testing.T) {
	files, e := embedAssets("./example/templates", "example/")
	failOnError(t, e)
	assert(t, len(files), 3)

	var (
		body []byte
		ok   bool
	)
	body, ok = files["templates/empty_file"]
	assert(t, ok, true)
	assert(t, len(body), 0)

	body, ok = files["templates/aaa.txt"]
	assert(t, ok, true)

	// Make sure it's gzipped data we can decompress again.
	buf := bytes.NewBuffer(body)
	gr, e := gzip.NewReader(buf)
	failOnError(t, e)
	defer gr.Close()
	tmp, e := ioutil.ReadAll(gr)
	failOnError(t, e)
	assert(t, string(tmp), "This file is just a test.\n")

}

func TestGenerateAssets(t *testing.T) {
	gen, e := NewGenerator("./example", "", "./example/assets/", "", "")
	failOnError(t, e)

	assert(t, gen.FilePath, "./example")
	assert(t, gen.Package, "assets")
	assert(t, gen.OutputPrefix, "./example/assets/")
	assert(t, gen.BuildTag, "embed")

	files, e := gen.GenerateAssets()
	failOnError(t, e)
	assert(t, len(files), 2)
	assert(t, files[0].Path, "./example/assets/dev.go")
	assert(t, files[1].Path, "./example/assets/build.go")
	if len(files[0].Body) < 1 {
		t.Error("Wasn't expecting an empty dev file")
	}
	if len(files[1].Body) < 1 {
		t.Error("Wasn't expecting an empty build file")
	}
}