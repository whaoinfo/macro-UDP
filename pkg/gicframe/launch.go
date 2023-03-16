package gicframe

import (
	"encoding/json"
	"github.com/whaoinfo/go-box/logger"
	"io/ioutil"
	"path"
)

type componentConfigModel struct {
	ComponentType string                 `json:"component_type"`
	Disable       bool                   `json:"disable"`
	Kw            map[string]interface{} `json:"kw"`
}

type configInfoModel struct {
	Key  string `json:"key"`
	Path string `json:"path"`
}

type LauncherConfigModel struct {
	AppID         string                     `json:"app_id"`
	LogLevel      string                     `json:"log_level"`
	ConfigInfoMap map[string]configInfoModel `json:"config_paths"`
	Components    []componentConfigModel     `json:"components"`
}

func LaunchDaemonApplication(workPath string, launcherConfPath string, newApp NewApplication,
	appArgs []interface{}, confInst IConfig, enabledDevMode bool) error {
	// Load launcher config
	if enabledDevMode {
		launcherConfPath = path.Join(workPath, "config_template", "launcher.json")
	}

	logger.InfoFmt("Launcher config, path: %s", launcherConfPath)
	fileData, readFileErr := ioutil.ReadFile(launcherConfPath)
	if readFileErr != nil {
		return readFileErr
	}
	launcherConf := &LauncherConfigModel{}
	if err := json.Unmarshal(fileData, launcherConf); err != nil {
		return err
	}

	// Set logger level
	logger.SetDefaultLogLevel(launcherConf.LogLevel)

	// Initialize APP
	logger.InfoFmt("Initialize application...")
	var app IApplication
	if newApp != nil {
		app = newApp()
	} else {
		app = &BaseApplication{}
	}
	if err := app.baseInitialize(ApplicationID(launcherConf.AppID)); err != nil {
		return err
	}
	if err := app.Initialize(appArgs...); err != nil {
		return err
	}
	logger.InfoFmt("The application has Initialized")

	// Initialize config
	logger.InfoFmt("Initialize config...")
	if err := initializeConfig(workPath, confInst, launcherConf, enabledDevMode); err != nil {
		return err
	}
	logger.Info("The config has initialized")

	setAppProxy(app)

	// Import component list
	logger.InfoFmt("Import components...")
	if err := app.importComponents(launcherConf.Components); err != nil {
		return err
	}
	logger.InfoFmt("The components have been imported")

	logger.Info("Start application...")
	if err := app.start(); err != nil {
		return err
	}
	logger.Info("The application has started")
	if err := app.AfterStart(); err != nil {
		return err
	}

	app.forever()

	if err := app.StopBefore(); err != nil {
		return err
	}
	logger.Info("Stop application...")
	if err := app.stop(); err != nil {
		return err
	}
	logger.Info("The application has stopped")

	return nil
}
