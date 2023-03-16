package gicframe

import (
	"encoding/json"
	"fmt"
	"github.com/whaoinfo/go-box/logger"
	"io/ioutil"
	"path"
)

type IConfig interface {
	Cover(data []byte) error
	Update(key string, data []byte) error
}

func initializeConfig(workPath string, confInst IConfig, launcherCfg *LauncherConfigModel, enabledDevMode bool) error {
	if confInst == nil {
		return nil
	}

	if enabledDevMode {
		for _, elem := range launcherCfg.ConfigInfoMap {
			elem.Path = path.Join(workPath, "config_template", fmt.Sprintf("%s.json", elem.Key))
		}
	}

	configDataMap := make(map[string]map[string]interface{})
	for _, info := range launcherCfg.ConfigInfoMap {
		logger.InfoFmt("Load config, key: %s, path: %s", info.Key, info.Path)
		fileData, readFileErr := ioutil.ReadFile(info.Path)
		if readFileErr != nil {
			return readFileErr
		}
		dataMap := make(map[string]interface{})
		if err := json.Unmarshal(fileData, &dataMap); err != nil {
			return err
		}
		configDataMap[info.Key] = dataMap
	}

	coverData, msErr := json.Marshal(configDataMap)
	if msErr != nil {
		return msErr
	}

	if err := confInst.Cover(coverData); err != nil {
		return err
	}

	return nil
}
