package yaber

import (
	"os"
	"testing"
)

func failOnError(t *testing.T, e error) {
	if e != nil {
		t.Errorf("Unexpected error: %v\n", e)
	}
}

func expectPkgName(t *testing.T, expected, path string) error {
	pkg, e := getPackageName(path)
	if e != nil {
		return e
	}
	if pkg != expected {
		t.Errorf("Bad package name, wanted: %s, got: %s\n", expected, pkg)
	}
	return nil
}

func TestPackageNameNotGo(t *testing.T) {
	e := expectPkgName(t, "templates", "./example/templates")
	failOnError(t, e)
}

func TestPackageNameIsGo(t *testing.T) {
	e := expectPkgName(t, "main", "./example")
	failOnError(t, e)
}

func TestPackageNameBadPath(t *testing.T) {
	e := expectPkgName(t, "", "./example/notadir")
	if _, ok := e.(*os.PathError); !ok {
		t.Errorf("Expected os.PathError, got: %#v\n", e)
	}
}
