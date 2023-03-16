package configmodel

import (
	"encoding/json"
	"github.com/whaoinfo/go-box/logger"
)

var (
	confInst = &Config{
		sidecarConfigModel: &SidecarConfigModel{},
		tracingModel:       &ConfigTracingModel{},
		simStorageModel:    &SimStorageModel{},
	}
)

type ConfigModel struct {
	SidecarConfigModel SidecarConfigModel `json:"sidecar"`
	ConfigTracingModel ConfigTracingModel `json:"tracing"`
	SimStorageModel    SimStorageModel    `json:"sim_storage"`
}

type Config struct {
	sidecarConfigModel *SidecarConfigModel
	tracingModel       *ConfigTracingModel
	simStorageModel    *SimStorageModel
}

func (t *Config) Cover(data []byte) error {
	md := &ConfigModel{}
	if err := json.Unmarshal(data, md); err != nil {
		return err
	}

	t.sidecarConfigModel = &md.SidecarConfigModel
	t.tracingModel = &md.ConfigTracingModel
	t.simStorageModel = &md.SimStorageModel

	logger.InfoFmt("config data: %v", string(data))
	return nil
}

func (t *Config) Update(key string, data []byte) error {
	// todo:

	return nil
}

func GetConfig() *Config {
	return confInst
}
