package main

import (
	"github.com/whaoinfo/go-box/logger"
	"github.com/whaoinfo/macro-UDP/gateway"
	"github.com/whaoinfo/macro-UDP/ioadapter"
	_ "github.com/whaoinfo/macro-UDP/message"
	"os"
	"time"
)

func newNetIO() (*gateway.NetIO, error) {
	plugs := []gateway.PlugInfo{
		{
			Type:                gateway.UDPPugType,
			ID:                  "UPD-SVR-33-00",
			Addr:                "0.0.0.0:7777",
			ReadQueueNum:        10,
			WriteQueueNum:       10,
			ReadQueueWorkerNum:  3,
			WriteQueueWorkerNum: 3,
			AllocateBufferFunc:  ioadapter.AllocateReadBufferFunc,
			ReadIOCallbackFunc:  ioadapter.OnHandleReadPacket,
		},
	}

	netIO := &gateway.NetIO{}
	if err := netIO.Initialize(plugs); err != nil {
		return nil, err
	}

	return netIO, nil
}

func main() {
	logger.SetDefaultLogLevel("ALL")
	netIO, newErr := newNetIO()
	if newErr != nil {
		logger.ErrorFmt("NewNetIO Failed, %v", newErr)
		os.Exit(1)
	}

	if err := netIO.Start(); err != nil {
		logger.ErrorFmt("NetIO has failed to start, %v", newErr)
		os.Exit(1)
	}

	for {
		time.Sleep(time.Second * 1)
	}
	os.Exit(0)
}
