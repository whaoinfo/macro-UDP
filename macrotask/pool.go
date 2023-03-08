package macrotask

import (
	"sync"
	"sync/atomic"
)

type Pool struct {
	workerNum     uint64
	busyWorkerNum uint64
	idleWorkerNum uint64

	rwMutex sync.RWMutex
}

func (t *Pool) OccupyWorker() (*Worker, error) {
	if atomic.AddUint64(&t.busyWorkerNum, 1) > t.workerNum {
		atomic.AddUint64(&t.busyWorkerNum, -1)
		return nil, nil
	}

	return nil, nil
}
