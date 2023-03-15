package storage

import (
	"errors"
	configmodel2 "github.com/whaoinfo/macro-UDP/internal/configmodel"
	"github.com/whaoinfo/macro-UDP/internal/define"
	frame "github.com/whaoinfo/macro-UDP/pkg/gicframe"
	"github.com/whaoinfo/macro-UDP/pkg/simpleworkerpool"
	sa "github.com/whaoinfo/macro-UDP/pkg/storageagent"
)

type ComponentKW struct {
	EnableStatsMode bool `json:"enable_stats_mode"`
	QueueGroup      struct {
		Maxsize   int64 `json:"maxsize"`
		MaxLength int64 `json:"maxlength"`
	} `json:"queue_group"`
	Worker struct {
		MaxSize         int64 `json:"maxsize"`
		IntervalMS      int   `json:"interval_ms"`
		IntervalRWCount int   `json:"interval_rw_count"`
	} `json:"worker"`
}

type Component struct {
	frame.BaseComponent
	agent      *sa.Agent
	workerPool *simpleworkerpool.Pool
	queueGroup *QueueGroup
	kw         *ComponentKW
}

func (t *Component) GetType() frame.ComponentType {
	return define.StorageComponentType
}

func (t *Component) Initialize(kw frame.IComponentKW) error {
	cfg := &configmodel2.GetConfigModel().Storage
	// initialize storage agent
	if err := t.initializeAgent(cfg); err != nil {
		return err
	}

	t.kw = kw.(*ComponentKW)
	// initialize queue group
	t.queueGroup = &QueueGroup{}
	t.queueGroup.Initialize(t.kw.QueueGroup.Maxsize, t.kw.QueueGroup.MaxLength)

	// initialize worker pool
	workerPool := simpleworkerpool.NewWorkerPool()
	if err := workerPool.Initialize(t.kw.Worker.MaxSize, t.kw.EnableStatsMode); err != nil {
		return err
	}
	t.workerPool = workerPool

	return nil
}

func (t *Component) initializeAgent(cfg *configmodel2.ConfigStorageModel) error {
	var infoList []sa.ClientInfo
	for _, tpy := range cfg.StorageAgent.ImportClientTypes {
		infoList = append(infoList, sa.ClientInfo{
			ClientType: sa.ClientType(tpy),
			Args:       nil,
		})
	}

	agent := &sa.Agent{}
	if err := agent.Initialize(infoList); err != nil {
		return err
	}

	t.agent = agent
	return nil
}

func (t *Component) Start() error {
	t.workerPool.Start()
	poolMaxsize := t.workerPool.GetMaxsize()
	for i := 0; i < int(poolMaxsize); i++ {
		if !t.workerPool.SubmitTask(t.timeHandleTask) {
			return errors.New("failed to submit task")
		}
	}

	return nil
}

//func (t *Component) PutElement(elem *ioadapter.LayerPacket) {
//
//}
