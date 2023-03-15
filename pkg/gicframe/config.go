package gicframe

import (
	"github.com/whaoinfo/go-box/logger"
	"io/ioutil"
	"path"
)

type IConfig interface {
	Parse(key string, data []byte) error
}

func initializeConfig(workPath string, confInst IConfig, launcherCfg *LauncherConfigModel, enabledDevMode bool) error {
	if confInst == nil {
		return nil
	}

	var infoList []configInfoModel
	if enabledDevMode {
		infoList = append(infoList, configInfoModel{
			Key: "default", Path: path.Join(workPath, "config_template", "config.json"),
		})
	} else {
		for _, elem := range launcherCfg.ConfigInfoMap {
			infoList = append(infoList, configInfoModel{
				Key: elem.Key, Path: elem.Path,
			})
		}
	}

	for _, info := range infoList {
		logger.InfoFmt("Load config, key: %v, path: %s", info.Key, info.Path)
		fileData, readFileErr := ioutil.ReadFile(info.Path)
		if readFileErr != nil {
			return readFileErr
		}
		if err := confInst.Parse(info.Key, fileData); err != nil {
			return err
		}
	}

	return nil
}
