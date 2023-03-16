package storage

import (
	"github.com/whaoinfo/go-box/logger"
	configmodel2 "github.com/whaoinfo/macro-UDP/internal/configmodel"
	"github.com/whaoinfo/macro-UDP/internal/define"
	frame "github.com/whaoinfo/macro-UDP/pkg/gicframe"
	"github.com/whaoinfo/macro-UDP/pkg/simpleworkerpool"
	sa "github.com/whaoinfo/macro-UDP/pkg/storageagent"
)

type ComponentKW struct {
	EnableStatsMode bool   `json:"enable_stats_mode"`
	AgentClientType string `json:"agent_client_type"`
	Worker          struct {
		MaxSize         int64 `json:"maxsize"`
		IntervalMS      int   `json:"interval_ms"`
		IntervalRWCount int   `json:"interval_rw_count"`
	} `json:"worker"`
}

type Component struct {
	frame.BaseComponent
	agent      *sa.Agent
	clientType sa.ClientType
	workerPool *simpleworkerpool.Pool
	kw         *ComponentKW
}

func (t *Component) GetType() frame.ComponentType {
	return define.StorageComponentType
}

func (t *Component) Initialize(kw frame.IComponentKW) error {
	cfg := &configmodel2.GetConfigModel().Storage
	// initialize storage agent
	kwArgs := kw.(*ComponentKW)
	t.clientType = sa.ClientType(kwArgs.AgentClientType)
	if err := t.initializeAgent(cfg); err != nil {
		return err
	}

	t.kw = kw.(*ComponentKW)
	// initialize worker pool
	workerPool := simpleworkerpool.NewWorkerPool()
	if err := workerPool.Initialize(t.kw.Worker.MaxSize, t.kw.EnableStatsMode); err != nil {
		return err
	}
	t.workerPool = workerPool

	return nil
}

func (t *Component) initializeAgent(cfg *configmodel2.ConfigStorageModel) error {
	var agentInfoList []sa.ClientInfo
	agentInfoList = append(agentInfoList, sa.ClientInfo{
		ClientType: sa.ClientType(tpy),
		Args:       []interface{}{cfg.StorageAgent.AmazonS3.Endpoint},
	})

	agent := &sa.Agent{}
	if err := agent.Initialize(agentInfoList); err != nil {
		return err
	}
	t.agent = agent

	if err := frame.GetAppProxy().Sub(define.StartTracSessionEvent, t.GetID(), t.onStartTracSessionEvent); err != nil {
		return err
	}

	return nil
}

func (t *Component) Start() error {
	t.workerPool.Start()
	return nil
}

func (t *Component) onStartTracSessionEvent(args ...interface{}) {
	if len(args) <= 0 {
		return
	}

	sessInfo := args[0].(*define.SessionStorageInfo)
	if !t.workerPool.SubmitTask(t.timeTask, args...) {
		logger.WarnFmt("Failed to submit a time task for %v session", sessInfo.ID)
		return
	}
	logger.AllFmt("The time task was submitted for %v session", sessInfo.ID)
}
