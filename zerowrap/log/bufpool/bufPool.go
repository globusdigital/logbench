package bufpool

import (
	"bytes"
	"sync"
)

// bufPool is a pool of *bytes.Buffer
type bufPool struct {
	pool *sync.Pool
}

// newBufPool creates a new bytes.Buffer pool
func newBufPool() bufPool {
	return bufPool{
		pool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

// Get returns a buffer
func (sbp bufPool) Get() *bytes.Buffer {
	return sbp.pool.Get().(*bytes.Buffer)
}

// Put puts the given buffer back to the pool
func (sbp bufPool) Put(buf *bytes.Buffer) {
	// @see https://go-review.googlesource.com/c/go/+/136116/4/src/fmt/print.go
	// Proper usage of a sync.Pool requires each entry to have approximately
	// the same memory cost. To obtain this property when the stored type
	// contains a variably-sized buffer, we add a hard limit on the
	// maximum buffer to place back in the pool.
	//
	// See https://golang.org/issue/23199
	const maxBufSize = 1 << 16 // 64KiB
	if buf.Cap() > maxBufSize {
		return
	}

	buf.Reset()
	sbp.pool.Put(buf)
}

// BufferPool is the global buffer pool used for error stack-trace generation
// to reduce the amount of allocated buffers
var BufferPool = newBufPool()
