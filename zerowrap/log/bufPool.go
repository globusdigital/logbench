package log

import (
	"bytes"
	"sync"
)

// bufPool is a bytes.Buffer pool
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
	buf.Truncate(0)
	sbp.pool.Put(buf)
}

var bufferPool = newBufPool()
