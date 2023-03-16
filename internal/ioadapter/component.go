package ioadapter

import (
	"fmt"
	"github.com/whaoinfo/go-box/logger"
	"github.com/whaoinfo/go-box/nbuffer"
	"github.com/whaoinfo/macro-UDP/internal/bufferpool"
	"github.com/whaoinfo/macro-UDP/internal/message"
	frame "github.com/whaoinfo/macro-UDP/pkg/gicframe"
	"github.com/whaoinfo/netio"
)

type ComponentKW struct {
	PluginType     string `json:"plugin_type"`
	Address        string `json:"address"`
	QueueNum       int    `json:"queue_num"`
	QueueWorkerNum int    `json:"queue_worker_num"`
}

type Component struct {
	frame.BaseComponent
	netIO *netio.NetIO
}

func (t *Component) Initialize(kw frame.IComponentKW) error {
	kwArgs := kw.(*ComponentKW)
	if err := t.initializeNetIO(kwArgs); err != nil {
		logger.WarnFmt("The NetIO package has failed to initialize, %v", err)
		return err
	}
	logger.InfoFmt("The NetIO package has initialized")
	return nil
}

func (t *Component) initializeNetIO(kw *ComponentKW) error {
	plugin := netio.PlugInfo{
		Type:                netio.PlugType(kw.PluginType),
		ID:                  netio.PlugID(fmt.Sprintf("%v_1", kw.PluginType)),
		Addr:                kw.Address,
		ReadQueueWorkerNum:  kw.QueueWorkerNum,
		WriteQueueWorkerNum: kw.QueueWorkerNum,
		ReadQueueNum:        kw.QueueNum,
		WriteQueueNum:       kw.QueueNum,
		AllocateBufferFunc:  t.allocateReadBufferFunc,
		ReadIOCallbackFunc:  t.handleReadPacket,
	}
	logger.InfoFmt("Initialize NetIO package, PluginType: %v, ID: %v, Address: %v, ReadQueueWorkerNum: %d, "+
		"WriteQueueWorkerNum: %d, ReadQueueNum: %d, WriteQueueNum: %d", plugin.Type, plugin.ID, plugin.Addr,
		plugin.ReadQueueWorkerNum, plugin.WriteQueueWorkerNum, plugin.ReadQueueNum, plugin.WriteQueueNum)
	netIO := &netio.NetIO{}
	if err := netIO.Initialize([]netio.PlugInfo{plugin}); err != nil {
		return err
	}

	t.netIO = netIO
	return nil
}

func (t *Component) Start() error {
	return t.netIO.Start()
}

func (t *Component) Stop() error {
	return t.netIO.Stop()
}

func (t *Component) allocateReadBufferFunc() (*nbuffer.BufferObject, error) {
	return bufferpool.AllocateBufferObject(), nil
}

func (t *Component) handleReadPacket(pkt *netio.UDPPacket) {
	if pkt == nil {
		logger.WarnFmt("The packet is a nil pointer")
		return
	}
	if pkt.Payload == nil {
		logger.WarnFmt("The payload of the packet is a nil value")
		return
	}

	buf := pkt.Payload.(*nbuffer.BufferObject)
	pkt.Payload = nil

	msgType, msg, regMsgInfo, ubErr := LayerPacketUnmarshalBinary(buf)
	if ubErr != nil {
		// todo
		logger.WarnFmt("Failed to call LayerPacket.UnmarshalBinary function, %v", ubErr)
		return
	}

	// put pkt
	if regMsgInfo.Handle == nil {
		// todo
		logger.WarnFmt("The RegMsgInfo.MessageHandle of message type %v is nil", msgType)
		return
	}

	ctx := &message.HandleContext{
		Buf: buf,
		Msg: msg,
	}
	regMsgInfo.Handle(ctx, regMsgInfo.Args...)

	return
}
