#!/bin/sh

# Example script for generating the asset files, inside the assets/ dir.

go run ../cmd/yaber/main.go -prefix="assets/" templates/
