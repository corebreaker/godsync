package sync

import (
	"runtime"
	"sync"
)

type tFileList struct {
	list []*FileDesc
	m    sync.Mutex
}

func (fl *tFileList) Len() int           { return len(fl.list) }
func (fl *tFileList) Less(i, j int) bool { return fl.list[i].Path < fl.list[j].Path }
func (fl *tFileList) Swap(i, j int)      { fl.list[i], fl.list[j] = fl.list[j], fl.list[i] }

func (fl *tFileList) add(fd *FileDesc) {
	fl.m.Lock()
	defer fl.m.Unlock()

	fl.list = append(fl.list, fd)
}

func (fl *tFileList) writer() chan<- *FileDesc {
	res := make(chan *FileDesc)
	for i := 0; i < runtime.GOMAXPROCS(-1); i++ {
		go func() {
			for fd := range res {
				fl.add(fd)
			}
		}()
	}

	return res
}
