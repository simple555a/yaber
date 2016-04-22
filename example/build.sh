#!/bin/sh

# Example script for generating the asset files, inside the assets/ dir.

cd ./assets/
go run ../../cmd/yaber.go --path "../templates/"
