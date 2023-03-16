package define

import (
	"github.com/whaoinfo/go-box/queue/safetyqueue"
	"github.com/whaoinfo/macro-UDP/pkg/gicframe"
)

const (
	IOAdapterComponentType gicframe.ComponentType = "ioadapter"
	TracerComponentType                           = "tracer"
	StorageComponentType                          = "storage"
)

type SessionID string

type GetSafetyQueueListFunc func() []*safetyqueue.SafetyQueue

type SessionStorageInfo struct {
	ID                     SessionID
	DisableStatus          *bool
	UploadCount            *uint64
	UploadFailCount        *uint64
	GetSafetyQueueListFunc GetSafetyQueueListFunc
}
