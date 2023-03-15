package storage

import (
	"github.com/whaoinfo/go-box/logger"
	"github.com/whaoinfo/macro-UDP/pkg/simpleworkerpool"
	"time"
)

func (t *Component) timeHandleTask(id simpleworkerpool.WorkerID, args ...interface{}) {
	waitDuration := time.Duration(t.kw.Worker.IntervalMS) * time.Millisecond
	queueID := int(id) % int(t.queueGroup.maxsize)

	for {
		time.Sleep(waitDuration)
		queue := t.queueGroup.GetQueueByID(queueID)
		queueID = (queueID + 1) % int(t.queueGroup.maxsize)
		if queue == nil {
			logger.WarnFmt("The queue is a nil pointer")
			continue
		}

		record := queue.Pops(t.kw.Worker.IntervalRWCount)
		//logger.AllFmt("Worker.TimeHandleTask, id: %v, pops len: %d", id, len(record))
		record = record

	}
}

func (t *Component) upload(record []interface{}) {
	//for _, d := range record {
	//
	//}
	//
	//t.agent.Upload(storageagent.AWSS3ClientType, "")
}
