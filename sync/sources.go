package sync

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type tSource struct {
	Path string // Absolute path of the source directory
	Name string //
	Root string // Name of the source directory used root directory in destination directory
}

func scanSource(src tSource, writer chan<- *FileDesc, done func()) {
	defer done()

	prefixSize := len([]rune(src.Path)) + 1
	batch := newHashingBatch()

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

		desc := &FileDesc{
			Path:  path,
			Name:  string([]rune(path)[prefixSize:]),
			Root:  src.Root,
			Size:  uint64(info.Size()),
			IsDir: false,
			Date:  info.ModTime(),
		}

		batch.add(desc)
		writer <- desc

		return nil
	}

	if err := filepath.Walk(src.Path, walker); err != nil {
		panic(err)
	}

	batch.wait()
}

func ScanSources(root string, dirs []string) []*FileDesc {
	registry := make(map[string]*tSource)
	counters := make(map[string]uint)

	for _, dir := range dirs {
		dirpath := dir
		dirname := filepath.Base(dirpath)

		idx := strings.IndexByte(dir, ':')
		if idx > 0 {
			dirpath = dir[(idx + 1):]
			dirname = dir[:idx]
		}

		if !filepath.IsAbs(dirpath) {
			dirpath = filepath.Clean(filepath.Join(root, dirpath))
		}

		prev := registry[dirname]
		if prev != nil {
			cnt, exists := counters[dirname]
			if exists {
				cnt++

				prev.Root += "_1"
				registry[prev.Root] = prev
			} else {
				cnt = 2
			}

			counters[dirname] = cnt
			dirname += fmt.Sprintf("_%d", cnt)
		}

		src := &tSource{
			Name: dirname,
			Root: dirname,
			Path: dirpath,
		}

		registry[src.Root] = src
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
