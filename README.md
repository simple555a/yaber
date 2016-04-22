yaber - Yet Another Binary Embedder
================================================================================

`yaber` is a yet another tool to generate go code with embedded binary data of
your assets.

- No external depencies
- Minimalistic and `gofmt`ed output
- Doesn't try to emulate fake files

Why
--------------------------------------------------------------------------------

As per Mars 2016, I had a hard time finding an existing go tool that would
let me embed my assets into my builds in a easy to use and non-buggy way.

Most of the other [projects](https://github.com/avelino/awesome-go#resource-embedding)
seemed to force you to import yet another depency just to let you load the
embedded files or failed to correctly implement the `http.File` interface without
introducing major bugs such as:

- Handles errors in bad ways
- [Infinite loops](https://github.com/GeertJohan/go.rice/issues/75)
- [Returns files from only the first dir](https://github.com/elazarl/go-bindata-assetfs/issues/14)

Ugh...

Quirks
--------------------------------------------------------------------------------

There is a couple of quirks though, with my own approach... No one's perfect,
I guess.

You will have to handle the embedded files as raw `[]byte` data, because we're
not trying to emulate `http.File`.

And this tool haven't been tested yet, so it's currently unknown how it will
handle large files or large amounts of files (is there a max size for a single
go source file? More research is needed).

Installation
--------------------------------------------------------------------------------

`go install github.com/lmas/yaber`

Usage
--------------------------------------------------------------------------------

Run:

`yaber -h`

to show the available options, or:

`yaber path/to/dir/with/assets`

It will generate two go files (as package `assets` by default, with the default
file prefix `asset_`) which contains the functions `Asset()` and `AssetDir()`.

By default you will load files directly from disk, using the `asset_dev.go` file.

If you enable the build flag `embed` (using `go run -tags "embed"` for example)
you will instead load embedded files found in the `asset_build.go` file.

`Asset(path string) ([]byte, error)`
Using this function, you can load files from `path` and get their binary data
as a `[]byte` slice.

`AssetDir(dir string) (map[string][]byte, error)`
This function will try to load all files in `path` and return a map where the
`string` keys are file paths and the `[]byte` values are the binary file data.
If an error is encountered, it will return a `nil` map and the `error`.

License
--------------------------------------------------------------------------------

MIT License, see the LICENSE file.

Todo
--------------------------------------------------------------------------------

- Prove that the code is non-buggy, by writting tests (when the project is nearing
  a mature state. Also not sur ehow to test the generated code...)
