package storage

import (
	"github.com/whaoinfo/go-box/logger"
	cfgmd "github.com/whaoinfo/macro-UDP/internal/configmodel"
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

	// initialize storage agent
	kwArgs := kw.(*ComponentKW)
	t.clientType = sa.ClientType(kwArgs.AgentClientType)
	if err := t.initializeAgent(); err != nil {
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

func (t *Component) initializeAgent() error {
	var agentInfo sa.ClientInfo
	switch t.clientType {
	case sa.AWSS3ClientType:
		s3Cfg := cfgmd.GetSidecarConfig().S3Storage
		agentInfo = sa.ClientInfo{
			ClientType: t.clientType,
			Args:       []interface{}{s3Cfg.Endpoint, "us-east-2", s3Cfg.AccessKeyID, s3Cfg.SecretAccessKey},
		}
		break
	case sa.SimClientType:
		simCfg := cfgmd.GetSimStorageConfig()
		agentInfo = sa.ClientInfo{
			ClientType: t.clientType,
			Args:       []interface{}{simCfg.Endpoint},
		}
		break
	}

	agent := &sa.Agent{}
	if err := agent.Initialize([]sa.ClientInfo{agentInfo}); err != nil {
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
