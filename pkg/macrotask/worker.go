package macrotask

import (
	"github.com/whaoinfo/go-box/logger"
	"runtime/debug"
	"time"
)

type WorkerHandle func(workerID int, args ...interface{})

type Worker struct {
	id           int
	intervalTime time.Duration
	handle       WorkerHandle
	handleArgs   []interface{}
	exitChan     chan bool
}

func (t *Worker) initialize(id int, intervalTime time.Duration, handle WorkerHandle, handleArgs ...interface{}) {
	t.id = id
	t.intervalTime = intervalTime
	t.handle = handle
	t.handleArgs = handleArgs
	t.exitChan = make(chan bool)
}

func (t *Worker) start() {
	go t.runForever()
	logger.AllFmt("ID %v worker has started", t.id)
}

func (t *Worker) stop() {
	logger.AllFmt("ID %v worker has stopped", t.id)
}

func (t *Worker) close() {
	close(t.exitChan)
}

func (t *Worker) runForever() {
	for {
		needBreak := false
		ticker := time.NewTicker(t.intervalTime)
		select {
		case <-ticker.C:
			t.executeHandleFunc()
			break
		case _, _ = <-t.exitChan:
			needBreak = true
			break
		}

		ticker.Stop()
		if needBreak {
			break
		}
	}

	logger.AllFmt("ID %v worker.runForever has broken", t.id)
	t.stop()
}

func (t *Worker) executeHandleFunc() {
	defer func() {
		if r := recover(); r != nil {
			logger.ErrorFmt("Catch the exception, recover: %v, stack: %v", r, string(debug.Stack()))
		}
	}()

	logger.AllFmt("Call worker handle, id: %v", t.id)
	t.handle(t.id, t.handleArgs...)
}

func NewWorkerPool() *WorkerPool {
	return &WorkerPool{}
}

type WorkerPool struct {
	poolSize  int
	workerMap map[int]*Worker
}

func (t *WorkerPool) Initialize(poolSize int, intervalTime time.Duration, handle WorkerHandle, handleArgs ...interface{}) {
	t.poolSize = poolSize
	t.workerMap = make(map[int]*Worker)

	for i := 0; i < poolSize; i++ {
		worker := &Worker{}
		worker.initialize(i, intervalTime, handle, handleArgs)
		t.workerMap[worker.id] = worker
	}
}

func (t *WorkerPool) Start() {
	for _, worker := range t.workerMap {
		worker.start()
	}
}

func (t *WorkerPool) Stop() {
	for _, worker := range t.workerMap {
		worker.close()
	}
}
