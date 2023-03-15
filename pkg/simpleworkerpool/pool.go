package simpleworkerpool

import (
	"errors"
	"github.com/whaoinfo/go-box/logger"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

func NewWorkerPool() *Pool {
	return &Pool{}
}

type Pool struct {
	maxsize          int64
	size             int64
	workerIncNum     int
	enableStatsMode  bool
	runningWorkerNum int64
	busyWorkerNum    int64
	stTaskCount      int64
	unStTaskCount    int64
	lock             sync.RWMutex
	workerMap        map[WorkerID]*Worker
}

func (t *Pool) Initialize(maxsize int64, enableStatsMode bool) error {
	if maxsize <= 0 {
		return errors.New("maxsize <= 0")
	}

	t.maxsize = maxsize
	t.workerMap = make(map[WorkerID]*Worker, maxsize)
	t.enableStatsMode = enableStatsMode

	return nil
}

func (t *Pool) GetMaxsize() int64 {
	return t.maxsize
}

func (t *Pool) Start() {
	for _, worker := range t.workerMap {
		worker.start()
	}
}

func (t *Pool) Stop() {
	for _, worker := range t.workerMap {
		worker.close()
	}
}

func (t *Pool) addRunningWorkerNum(num int64) {
	atomic.AddInt64(&t.runningWorkerNum, num)
}

func (t *Pool) addBusyWorkerNum(num int64) {
	atomic.AddInt64(&t.busyWorkerNum, num)
}

func (t *Pool) SubmitTask(taskFunc WorkerHandle, args ...interface{}) (submitted bool) {
	t.lock.Lock()
	defer func() {
		if r := recover(); r != nil {
			logger.ErrorFmt("Catch the exception, recover: %v, stack: %v", r, string(debug.Stack()))
		}

		if !submitted {
			t.unStTaskCount += 1
		}
	}()

	t.stTaskCount += 1
	if t.size >= t.maxsize {
		t.lock.Unlock()
		return
	}

	worker := t.getIdleWorker()
	if worker == nil {
		worker = t.addWorker()
	}
	t.lock.Unlock()

	ctxt := &Context{
		Event:  ExecuteEvent,
		Handle: taskFunc,
		Args:   args,
	}
	worker.pub(ctxt)

	submitted = true
	return
}

func (t *Pool) generateWorkerID() WorkerID {
	t.workerIncNum += 1
	return WorkerID(t.workerIncNum)
}

func (t *Pool) addWorker() *Worker {
	worker := &Worker{}
	worker.initialize(t.generateWorkerID(), t)
	t.workerMap[worker.GetID()] = worker
	t.size += 1
	worker.start()
	logger.AllFmt("Add a worker, id: %v", worker.GetID())
	return worker
}

func (t *Pool) getIdleWorker() *Worker {

	// todo
	return nil
}

type PoolRangeFunc func(id WorkerID, startTimestamp, executeTimestamp int64, status StatusType)

func (t *Pool) Range(f PoolRangeFunc) {
	t.lock.RLock()
	workers := make([]*Worker, 0, len(t.workerMap))
	for _, worker := range t.workerMap {
		workers = append(workers, worker)
	}
	t.lock.RUnlock()

	for _, worker := range workers {
		f(worker.GetID(), worker.GetStartTimestamp(), worker.GetExecuteTimestamp(), worker.GetStatus())
	}

}
