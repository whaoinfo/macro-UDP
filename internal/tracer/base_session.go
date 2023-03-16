package tracer

import (
	"errors"
	"github.com/whaoinfo/go-box/ctime"
	"github.com/whaoinfo/go-box/logger"
	"github.com/whaoinfo/go-box/queue/ringqueue"
	"github.com/whaoinfo/go-box/queue/safetyqueue"
	"github.com/whaoinfo/macro-UDP/internal/configmodel"
	"github.com/whaoinfo/macro-UDP/internal/define"
	message "github.com/whaoinfo/macro-UDP/internal/message"
	frame "github.com/whaoinfo/macro-UDP/pkg/gicframe"
	"sync/atomic"
)

func NewBaseSession() ISession {
	return &BaseSession{}
}

type BaseSession struct {
	id           define.SessionID
	disable      bool
	endTimestamp int64
	filters      []IFilter
	queueList    []*safetyqueue.SafetyQueue
	putCount     int64
	putFailCount int64

	uploadCount     uint64
	uploadFailCount uint64

	newFilters NewFiltersFunc
}

func (t *BaseSession) initialize(maxQueueNum, queueLength int, cfg *configmodel.ConfigTraceSessionModel) error {
	if cfg.TraceSessionId == "" {
		return errors.New("the traceSessionId field in the configuration is empty")
	}
	if cfg.EndTime == "" {
		return errors.New("the endTime field in the configuration is empty")
	}

	t.id = define.SessionID(cfg.TraceSessionId)
	t.disable = cfg.Disable
	for i := 0; i < maxQueueNum; i++ {
		q := safetyqueue.NewSafetyQueue(queueLength, func(length int) safetyqueue.IQueue {
			return ringqueue.NewRingQueue(length)
		})
		t.queueList = append(t.queueList, q)
	}

	filters, err := t.newFilters(cfg)
	if err != nil {
		return err
	}

	t.filters = append(t.filters, filters...)
	return nil
}

func (t *BaseSession) getID() define.SessionID {
	return t.id
}

func (t *BaseSession) getQueueList() []*safetyqueue.SafetyQueue {
	return t.queueList
}

func (t *BaseSession) start() error {
	info := &define.SessionStorageInfo{
		ID:                     t.getID(),
		DisableStatus:          &t.disable,
		UploadCount:            &t.uploadCount,
		UploadFailCount:        &t.uploadFailCount,
		GetSafetyQueueListFunc: t.getQueueList,
	}
	if err := frame.GetAppProxy().Pub(define.StartTracSessionEvent, info); err != nil {
		return err
	}

	return nil
}

func (t *BaseSession) stop() error {

	return nil
}

func (t *BaseSession) match(msg message.IMessage) bool {
	if t.disable {
		return false
	}
	if t.endTimestamp >= ctime.CurrentTimestamp() {
		return false
	}

	for _, filter := range t.filters {
		if !filter.Check(msg) {
			continue
		}
		return true
	}

	return false
}

func (t *BaseSession) putMessageContext(ctx *message.HandleContext) bool {
	c := atomic.AddInt64(&t.putCount, 1)
	idx := int(c) % len(t.queueList)
	q := t.queueList[idx]
	if !q.Put(ctx) {
		atomic.AddInt64(&t.putFailCount, 1)
		return false
	}

	logger.AllFmt("The context was put in the queue with index %d", idx)
	return true
}
