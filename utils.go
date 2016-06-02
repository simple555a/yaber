package yaber

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"
)

// Get the executed command and flags.
func executedCommand() string {
	return fmt.Sprintf("yaber %s", strings.Join(os.Args[1:], " "))
}

// Auto detect the GO pkg name for a file path.
func getPackageName(path string) (string, error) {
	dir, e := filepath.Abs(path)
	if e != nil {
		return "", e
	}
	// Use the dir as the default pkg name
	pkgName := filepath.Base(dir)

	// Check if the dir contains a go pkg
	pkg, e := build.ImportDir(dir, 0)
	if e != nil {
		// NoGoError = not a valid pkg in dir, so ignore it and use the
		// default from above.
		if _, ok := e.(*build.NoGoError); !ok {
			return "", e
		}
	} else {
		// Looks like it was an actual go pkg
		pkgName = pkg.Name
	}
	return pkgName, nil
}
