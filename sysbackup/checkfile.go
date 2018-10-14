package main

import (
	"os"

	"github.com/siddontang/go/ioutil2"
)

func checkdir(kind, path string) {
	if !ioutil2.FileExists(path) {
		panic(fmt.Errorf("The %s `%s` was not found", kind, path))
	}

	info, err := os.Stat(path)
	if err != nil {
		panic(err)
	}

	if !info.IsDir() {
		panic(fmt.Errorf("The %s `%s` is not a directory", kind, path))
	}
}
