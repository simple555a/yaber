// +build !embed

package assets

import "io/ioutil"

func Asset(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}
