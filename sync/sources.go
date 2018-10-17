package sync

import (
	"crypto/sha256"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

type tSource struct {
	Path string
	Name string
}

func scanDir(src tSource, writer chan<- *tFileList) {
	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		dest_files[path] = &info

		return nil
	}

	err := filepath.Walk(filepath.Join(dest, filepath.Base(src)), walker)

	if err != nil {
		panic(err)
	}
}

//dir_to_scan, dest
// dir_to_scan <= getRoot(root, dirpath)
func scanSource(src tSource, writer chan<- *tFileList, done func()) {
	defer done()
}

func ScanSources(root string, dirs []string) []*FileDesc {
    registry := make(map[string]*tSource)

	for _, dir := range dirs {
		dirpath := dir
		dirname := filepath.Base(dirpath)

		if !filepath.IsAbs(dir) {
			dirpath = filepath.Clean(filepath.Join(root, dirpath))
		}

        prev := registry[dirname]
        if prev != nil {

        }

        src := &tSource{
            Name: dirname,
            Path: dirpath,
        }


	}

	var wg sync.WaitGroup

	list := new(tFileList)
	wr := list.writer()
	defer close(wr)

	for _, dir := range dirs {
		wg.Add(1)
		go scanSource(root, dir, wr, wg.Done)
	}

	wg.Wait()
	sort.Sort(list)

	return list.list
}
