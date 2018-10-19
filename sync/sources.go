package sync

import (
	"crypto/sha256"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

type tSource struct {
	Path string // Absolute path of the source directory
	Name string //
	Root string // Name of the source directory used root directory in destination directory
}

func scanSource(src tSource, writer chan<- *FileDesc, done func()) {
	defer done()

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			writer <- &FileDesc{
				Path:  path,
				Root:  src.Root,
				IsDir: true,
			}

			return nil
		}

		return nil
	}

	err := filepath.Walk(src.Path, walker)

	if err != nil {
		panic(err)
	}
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

	for _, src := range registry {
		wg.Add(1)
		go scanSource(src, wr, wg.Done)
	}

	wg.Wait()
	sort.Sort(list)

	return list.list
}
