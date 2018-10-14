package main

import (
	"fmt"
	"os"
	"path/filepath"
    "crypto/sha256"
)

func doBackup(dest string, sources []string) {
	dest_files := make(map[string]*os.FileInfo)
	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		dest_files[path] = &info

		return nil
	}

	for _, src := range sources {
		err := filepath.Walk(filepath.Join(dest, filepath.Base(src)), walker)

		if err != nil {
			panic(err)
		}
	}

    sha256.Sum256()
}
