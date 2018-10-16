package sync

import "time"

type FileDesc struct {
	Path string
	Name string
	Hash string
	Size uint64
	Date *time.Time
}
