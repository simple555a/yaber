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

func TestEmbedAssetEmptyDir(t *testing.T) {
	files, e := embedAssets("./test/empty", "")
	failOnError(t, e)
	assert(t, len(files), 0)
}

func TestEmbedAssets(t *testing.T) {
	files, e := embedAssets("./test/files", "test/")
	failOnError(t, e)
	assert(t, len(files), 2)

	var (
		body []byte
		ok   bool
	)
	body, ok = files["files/empty.txt"]
	assert(t, ok, true)
	assert(t, len(body), 0)

	body, ok = files["files/notempty.txt"]
	assert(t, ok, true)

	// Make sure it's gzipped data we can decompress again.
	buf := bytes.NewBuffer(body)
	gr, e := gzip.NewReader(buf)
	failOnError(t, e)
	defer gr.Close()
	tmp, e := ioutil.ReadAll(gr)
	failOnError(t, e)
	assert(t, string(tmp), "here's a line in a file\n")

}

func TestGenerateAssets(t *testing.T) {
	gen, e := NewGenerator("./test/pkg", "", "./test/pkg/", "", "")
	failOnError(t, e)

	assert(t, gen.FilePath, "./test/pkg")
	assert(t, gen.Package, "main")
	assert(t, gen.OutputPrefix, "./test/pkg/")
	assert(t, gen.BuildTag, "embed")

	files, e := gen.GenerateAssets()
	failOnError(t, e)
	assert(t, len(files), 2)
	assert(t, files[0].Path, "./test/pkg/dev.go")
	assert(t, files[1].Path, "./test/pkg/build.go")
	if len(files[0].Body) < 1 {
		t.Error("Wasn't expecting an empty dev file")
	}
	if len(files[1].Body) < 1 {
		t.Error("Wasn't expecting an empty build file")
	}
}
