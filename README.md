yaber - Yet Another Binary Embedder
================================================================================

[![GoDoc](https://godoc.org/github.com/lmas/yaber?status.svg)](https://godoc.org/github.com/lmas/yaber)
[![Build Status](https://travis-ci.org/lmas/yaber.svg?branch=master)](https://travis-ci.org/lmas/yaber)
[![Coverage Status](https://coveralls.io/repos/github/lmas/yaber/badge.svg?branch=master)](https://coveralls.io/github/lmas/yaber?branch=master)

`yaber` is a yet another Go tool to generate code with embedded binary data,
from your assets.

Main features:

- Minimalistic and `gofmt`ed output.
- No external depencies for the final assets.
- Doesn't try to emulate fake files (`[]byte` slices is good enough).
- Try really hard to find good default values.
- Generates a companion test file.

Status
--------------------------------------------------------------------------------

Currently in beta, main functionality has been implemented and is currently
being polished and tested in production.

Installation
--------------------------------------------------------------------------------

Run

    go install github.com/lmas/yaber/cmd/yaber

to install the simple command line tool, or


    go install github.com/lmas/yaber

to just install the library.

Usage
--------------------------------------------------------------------------------

Run

    yaber -h

to show the available options, or

    yaber path/to/dir/with/assets

It will generate two Go files with the default names `assets.go` and
`assets_test.go`.

`assets.go` contains the three main functions:

`asset(path string) ([]byte, error)`
Using this function, you can load an embedded file from `path` and get it's
binary data as a `[]byte` slice.

`assetDir(dir string) (map[string][]byte, error)`
This function will try to load all files in `dir` and return a map where the
`string` keys are file paths and the `[]byte` values are the binary file data.

`setRawAssets(b bool)`
If set to true, `asset` and `assetDir` will load files directly from disk
instead.

Why
--------------------------------------------------------------------------------

As per Mars 2016, I had a hard time finding an existing go tool that would
let me embed my assets into my builds in a easy to use and non-buggy way.

Most of the other [projects](https://github.com/avelino/awesome-go#resource-embedding)
seemed to force you to import yet another depency, just to let you load the
embedded files, or failed to correctly implement the `http.File` interface
(or similar things), without introducing bugs and quirks.

Since then, I've not really bothered to find a tool that would meet my
expectations and instead focused on making this tool.

License
--------------------------------------------------------------------------------

MIT License, see the LICENSE file.

Todo
--------------------------------------------------------------------------------

- Need to improve the tests (more test cases and test errors).

- Is there a max size for a single go source file? Perform a more extensive test,
  using really large assets/large amounts of assets.

