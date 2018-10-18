package sync

import "time"

type FileDesc struct {
	Path  string
	Name  string
	Root  string
	Hash  string
	Size  uint64
	IsDir bool
	Date  *time.Time
}
