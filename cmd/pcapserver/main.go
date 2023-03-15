package main

import (
	"flag"
	"github.com/whaoinfo/go-box/logger"
	"github.com/whaoinfo/macro-UDP/internal/configmodel"
	_ "github.com/whaoinfo/macro-UDP/internal/ioadapter"
	_ "github.com/whaoinfo/macro-UDP/internal/storage"
	_ "github.com/whaoinfo/macro-UDP/internal/tracer"
	frame "github.com/whaoinfo/macro-UDP/pkg/gicframe"
	"os"
)

func main() {
	lachConfPath := flag.String("launcher_cfg_path", "", "launcher_cfg_path=")
	enableDevMode := flag.Bool("enable_dev_mode", false, "enable_dev_mode=false, true")
	flag.Parse()

	workPath, getWdErr := os.Getwd()
	if getWdErr != nil {
		logger.ErrorFmt("OS.Getwd Failed, %v", getWdErr)
		os.Exit(1)
	}

	cfgInst := configmodel.GetConfigInstance()
	if err := frame.LaunchDaemonApplication(workPath, *lachConfPath, nil, nil, cfgInst,
		*enableDevMode); err != nil {
		logger.ErrorFmt("Failed to launch application, %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}
