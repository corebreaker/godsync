package sync

import (
	"crypto/sha512"
	"encoding/hex"
	"io"
	"os"
	"sync"
)

const HASHING_POOL_SIZE = 16

var poolChan = make(chan func())

func init() {
	for i := 0; i < HASHING_POOL_SIZE; i++ {
		go func() {
			for f := range poolChan {
				f()
			}
		}()
	}
}

type tHashingBatch struct {
	wg sync.WaitGroup
}

func (hb *tHashingBatch) add(desc *FileDesc) {
	hb.wg.Add(1)
	poolChan <- func() {
		defer hb.wg.Done()

		f, err := os.Open(desc.Path)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		hasher := sha512.New()
		if _, err := io.Copy(hasher, f); err != nil {
			panic(err)
		}

		desc.Hash = hex.EncodeToString(hasher.Sum(nil))
	}
}

func (hb *tHashingBatch) wait() {
	hb.wg.Wait()
}

func newHashingBatch() *tHashingBatch {
	return new(tHashingBatch)
}
