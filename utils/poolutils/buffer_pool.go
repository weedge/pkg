package poolutils

// Cache pool to avoid repeated memory requests and reduce gc pressure
// A sync.Pool object is a collection of temporary objects. Pool is corout-safe
// Note: After getting the resource Get, you must perform the Put operation to put back

import (
	"bytes"
	"sync"
)

type BufferPool struct {
	sync.Pool
}

// NewBufferPool with buffer size
func NewBufferPool(bufferSize int) (bufferpool *BufferPool) {
	return &BufferPool{
		sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 0, bufferSize))
			},
		},
	}
}

// Get byte.Buffer from buff pool (sync.Pool)
func (bufferpool *BufferPool) Get() *bytes.Buffer {
	return bufferpool.Pool.Get().(*bytes.Buffer)
}

// Put byte.Buffer to buff pool (sync.Pool)
func (bufferpool *BufferPool) Put(b *bytes.Buffer) {
	b.Reset()
	bufferpool.Pool.Put(b)
}
