package ioadapter

import (
	"github.com/whaoinfo/go-box/logger"
	"github.com/whaoinfo/go-box/nbuffer"
	"github.com/whaoinfo/macro-UDP/bufferpool"
	"github.com/whaoinfo/macro-UDP/gateway"
)

func AllocateReadBufferFunc() (*nbuffer.BufferObject, error) {
	return bufferpool.AllocateBufferObject(), nil
}

func OnHandleReadPacket(pkt *gateway.UDPPacket) {
	if pkt == nil {
		logger.WarnFmt("The pkt is a nil pointer")
		return
	}
	if pkt.Payload == nil {
		logger.WarnFmt("The payload of the pkt is a nil value")
		return
	}

	buf := pkt.Payload.(*nbuffer.BufferObject)
	pkt.Payload = nil

	layerPkt := NewLayerPacket()
	if err := layerPkt.UnmarshalBinary(buf); err != nil {
		logger.WarnFmt("Failed to call LayerPacket.UnmarshalBinary function, %v", err)
	}

	// todo put layer pkt to macrotask

	return
}
