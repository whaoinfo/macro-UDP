package simpleworkerpool

import (
	"github.com/whaoinfo/go-box/ctime"
	"github.com/whaoinfo/go-box/logger"
	"runtime/debug"
)

type StatusType int

type WorkerID int

type WorkerHandle func(id WorkerID, args ...interface{})

type Worker struct {
	id               WorkerID
	pool             *Pool
	startTimestamp   int64
	executeTimestamp int64
	status           StatusType
	ctxtChan         chan *Context
	currentEvent     EventType
}

func (t *Worker) initialize(id WorkerID, pool *Pool) {
	t.id = id
	t.pool = pool
	t.status = MinNonState
	t.ctxtChan = make(chan *Context)
}

func (t *Worker) GetID() WorkerID {
	return t.id
}

func (t *Worker) GetStartTimestamp() int64 {
	return t.startTimestamp
}

func (t *Worker) GetExecuteTimestamp() int64 {
	return t.executeTimestamp
}

func (t *Worker) GetStatus() StatusType {
	return t.status
}

func (t *Worker) start() {
	t.startTimestamp = ctime.CurrentTimestamp()
	t.status = RunningStatus
	t.pool.addRunningWorkerNum(1)
	go t.runForever()
}

func (t *Worker) close() {
	close(t.ctxtChan)
}

func (t *Worker) runForever() {
	var ctxt *Context
	var ok bool
	for {
		t.status = IdleStatus
		ctxt, ok = <-t.ctxtChan
		if !ok {
			t.currentEvent = CloseEvent
			t.status = ClosedStatus
			logger.AllFmt("Worker channel %v has closed")
			break
		}

		event := ctxt.Event
		if ctxt.Event < MinNonEvent || ctxt.Event >= MaxNonEvent {
			logger.WarnFmt("Event type %v dose not exist", event)
			continue
		}

		t.currentEvent = event
		switch event {
		case ExecuteEvent:
			t.execute(ctxt)
			break
		default:
			logger.WarnFmt("Event type %v was not processed", event)
		}
	}

	t.stop(ctxt)
}

func (t *Worker) stop(ctxt *Context) {
	if ctxt != nil {

	}

	t.status = StoppedStatus
	logger.AllFmt("Worker %v has stopped", t.id)
}

func (t *Worker) pub(ctxt *Context) {
	t.ctxtChan <- ctxt
}

func (t *Worker) execute(ctxt *Context) {
	if ctxt.Handle == nil {
		logger.WarnFmt("Context.Handle is a nil pointer")
		return
	}

	t.pool.addBusyWorkerNum(1)
	defer func() {
		if r := recover(); r != nil {
			logger.ErrorFmt("Catch the exception, recover: %v, stack: %v", r, string(debug.Stack()))
		}

		t.currentEvent = ExceptionEvent
		t.status = IdleStatus
		t.pool.addBusyWorkerNum(-1)
	}()

	t.status = BusyStatus
	t.executeTimestamp = ctime.CurrentTimestamp()

	logger.AllFmt("Call worker handle, id: %v", t.id)
	ctxt.Handle(t.id, ctxt.Args...)
	logger.AllFmt("Called worker handle, id: %v", t.id)

	t.currentEvent = CompleteEvent

	return
}
