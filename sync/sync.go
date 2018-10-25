package sync

import "time"

type FileDesc struct {
	Path  string    // Absolute path of the source file
	Name  string    // Relative path to the root (common part to source and destination)
	Root  string    // Root directory name to be used in destination directory
	Hash  string    // SHA-256 hash of the source file
	Size  uint64    // Size of the source file
	IsDir bool      // The source file is a subdirectory ?
	Date  time.Time // Last modification time of the source file
}
