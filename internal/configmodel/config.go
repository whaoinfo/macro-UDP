package configmodel

import "encoding/json"

var (
	confModel = &ConfigModel{}
	confInst  = &Config{
		model: confModel,
	}
)

type ConfigModel struct {
	Tracer  ConfigTracerModel  `json:"tracer"`
	Storage ConfigStorageModel `json:"storage"`
}

type Config struct {
	model *ConfigModel
}

func (t *Config) Parse(key string, data []byte) error {
	if key == "default" {
		return t.parseDefault(data)
	}

	return nil
}

func (t *Config) parseDefault(data []byte) error {
	return json.Unmarshal(data, t.model)
}

func GetConfigInstance() *Config {
	return confInst
}

func GetConfigModel() *ConfigModel {
	return confInst.model
}
