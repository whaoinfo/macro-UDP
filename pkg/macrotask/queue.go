package macrotask

import (
	"github.com/whaoinfo/go-box/queue/ringqueue"
	"github.com/whaoinfo/go-box/queue/safetyqueue"
)

type QueuePool struct {
	poolSize    int
	queueLength int
	queueMap    map[int]*safetyqueue.SafetyQueue
}

func (t *QueuePool) Initialize(poolSize, queueLength int) {
	t.poolSize = poolSize
	t.queueLength = queueLength
	t.queueMap = make(map[int]*safetyqueue.SafetyQueue)
	for i := 0; i < poolSize; i++ {
		q := safetyqueue.NewSafetyQueue(queueLength, func(maxLength int) safetyqueue.IQueue {
			return ringqueue.NewRingQueue(queueLength)
		})
		t.queueMap[i] = q
	}
}

func (t *QueuePool) Puts(items ...interface{}) int {
	q := t.GetBalanceQueue()
	return q.Puts(items...)
}

func (t *QueuePool) Pops(count int) []interface{} {
	queue := t.GetBalanceQueue()
	return queue.Pops(count)
}

func (t *QueuePool) GetBalanceQueue() *safetyqueue.SafetyQueue {
	// todo
	id := 0
	queue := t.queueMap[id]
	return queue
}
