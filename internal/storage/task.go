package storage

import (
	"bytes"
	"github.com/whaoinfo/go-box/logger"
	"github.com/whaoinfo/go-box/queue/safetyqueue"
	"github.com/whaoinfo/macro-UDP/internal/define"
	"github.com/whaoinfo/macro-UDP/internal/message"
	"github.com/whaoinfo/macro-UDP/pkg/simpleworkerpool"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

const (
	workerBufBytesSize = 1024 * 1024 * 5
)

func (t *Component) timeTask(id simpleworkerpool.WorkerID, args ...interface{}) {
	sessInfo := args[0].(*define.SessionStorageInfo)
	logger.AllFmt("The time task of %v session has started, WorkerID: %v", sessInfo.ID, id)
	queueList := sessInfo.GetSafetyQueueListFunc()
	var queueIdx int
	if len(queueList) <= 0 {
		logger.WarnFmt("%v session queue number is 0, exit this time task", sessInfo.ID)
		return
	}

	uploadBuf := bytes.NewBuffer(make([]byte, workerBufBytesSize))
	queueIdx = int(id) % len(queueList)
	waitDuration := time.Duration(t.kw.Worker.IntervalMS) * time.Millisecond

	for {
		time.Sleep(waitDuration)
		if *sessInfo.DisableStatus {
			break
		}

		queueList = sessInfo.GetSafetyQueueListFunc()
		if len(queueList) <= 0 {
			logger.WarnFmt("%v session queue number is 0", sessInfo.ID)
			continue
		}

		elems := t.fetchElemsList(&queueIdx, queueList)
		if len(elems) <= 0 {
			continue
		}

		uploadBuf.Reset()
		if err := t.uploadElements(uploadBuf, sessInfo, elems); err != nil {
			logger.WarnFmt("Upload failed, %v", err)
		}

		// todo: recycle buffers
	}

	logger.AllFmt("Exit TimeTask, WorkerID: %v, SessionID: %V", id, sessInfo.ID)
}

func (t *Component) fetchElemsList(queueIdx *int, queueList []*safetyqueue.SafetyQueue) []interface{} {
	var retList []interface{}
	readNum := t.kw.Worker.IntervalRWCount
	queueListLen := len(queueList)

	for i := 0; i < queueListLen; i++ {
		q := queueList[*queueIdx]
		*queueIdx = (*queueIdx + 1) % queueListLen
		elems := q.Pops(readNum)
		logger.AllFmt("Fetch context on queue %d, len: %d", *queueIdx, len(elems))
		if len(elems) <= 0 {
			continue
		}

		retList = append(retList, elems...)
		readNum -= len(elems)
		if readNum <= 0 {
			break
		}
	}

	return retList
}

func (t *Component) uploadElements(uploadBuf *bytes.Buffer, sessInfo *define.SessionStorageInfo, elems []interface{}) error {
	uploadCount := atomic.AddUint64(sessInfo.UploadCount, 1)
	bucket := strings.Join([]string{string(sessInfo.ID), "0"}, "-")
	key := strconv.Itoa(int(uploadCount))

	for _, elem := range elems {
		ctx := elem.(*message.HandleContext)
		uploadBuf.Write(ctx.Buf.GetNextReadBytes())
	}

	if err := t.agent.UploadLarge(t.clientType, bucket, key, uploadBuf, workerBufBytesSize); err != nil {
		atomic.AddUint64(sessInfo.UploadFailCount, 1)
		return err
	}

	return nil
}
