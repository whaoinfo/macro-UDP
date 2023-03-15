package storage

import (
	"github.com/whaoinfo/go-box/queue/ringqueue"
	"github.com/whaoinfo/go-box/queue/safetyqueue"
)

type QueueGroup struct {
	maxsize        int64
	maxQueueLength int64
	queueMap       map[int]*safetyqueue.SafetyQueue
}

func (t *QueueGroup) Initialize(maxsize, maxQueueLength int64) {
	t.maxsize = maxsize
	t.maxQueueLength = maxQueueLength
	t.queueMap = make(map[int]*safetyqueue.SafetyQueue)
	maxsizeInt := int(maxsize)
	for i := 0; i < maxsizeInt; i++ {
		q := safetyqueue.NewSafetyQueue(maxsizeInt, func(maxLength int) safetyqueue.IQueue {
			return ringqueue.NewRingQueue(maxsizeInt)
		})
		t.queueMap[i] = q
	}
}

func (t *QueueGroup) Puts(items ...interface{}) int {
	q := t.GetBalanceQueue()
	return q.Puts(items...)
}

func (t *QueueGroup) Pops(count int) []interface{} {
	queue := t.GetBalanceQueue()
	return queue.Pops(count)
}

func (t *QueueGroup) GetBalanceQueue() *safetyqueue.SafetyQueue {
	// todo
	id := 0
	queue := t.queueMap[id]
	return queue
}

func (t *QueueGroup) GetQueueByID(queueID int) *safetyqueue.SafetyQueue {
	return t.queueMap[queueID]
}
