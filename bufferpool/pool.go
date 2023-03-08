package bufferpool

import (
	"github.com/whaoinfo/go-box/nbuffer"
	"sync"
)

const (
	defaultBuffCapacity = 1024 * 2
)

var (
	pool = &sync.Pool{
		New: func() any {
			return nbuffer.NewBufferObject(defaultBuffCapacity)
		},
	}
)

func AllocateBufferObject() *nbuffer.BufferObject {
	obj := pool.Get()
	return obj.(*nbuffer.BufferObject)
}

func RecycleBufferObject(obj *nbuffer.BufferObject) {
	if obj == nil {
		return
	}

	pool.Put(obj)
}

func ReleaseBufferObject(obj *nbuffer.BufferObject) {
	if obj == nil {
		return
	}
	obj.Release()
}
