package sync

import (
	"runtime"
	"sync"
)

type tFileList struct {
	list []*FileDesc
	m    sync.Mutex
}

func (fl *tFileList) add(fd *FileDesc) {
	fl.m.Lock()
	defer fl.m.Unlock()

	fl.list = append(fl.list, fd)
}

func (fm *tFileList) writer() chan<- *FileDesc {
	res := make(chan *FileDesc)
	for i := 0; i < runtime.GOMAXPROCS(-1); i++ {
		go func() {
			for fd := range res {
				fm.add(fd)
			}
		}()
	}

	return res
}
