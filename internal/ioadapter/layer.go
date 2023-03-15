package ioadapter

import (
	"encoding/binary"
	"fmt"
	"github.com/whaoinfo/go-box/logger"
	"github.com/whaoinfo/go-box/nbuffer"
	"github.com/whaoinfo/macro-UDP/internal/message"
)

const (
	MsgTypeSize   = 1
	MsgLengthSize = 2
	LayerHdrSize  = MsgTypeSize + MsgLengthSize
)

type LayerPacketHeader struct {
	MsgType   uint8
	MsgLength uint16
}

func LayerPacketUnmarshalBinary(buf *nbuffer.BufferObject) (message.MsgType, message.IMessage, *message.RegMessageInfo, error) {
	bufBytes := buf.Bytes()
	layerHdrBuf := bufBytes[0:LayerHdrSize]
	// parse message type
	msgType := message.MsgType(layerHdrBuf[0])
	logger.DebugFmt("Message type %v UnmarshalBinary", msgType)

	// parse message length
	msgLen := int(binary.BigEndian.Uint16(layerHdrBuf[1:]))
	// Check length
	if (buf.GetWriteLength() - LayerHdrSize) < msgLen {
		return msgType, nil, nil, fmt.Errorf("message type %v length does not match", msgType)
	}

	// Find message type
	regInfo := message.GetRegisterMessageInfo(msgType)
	if regInfo == nil {
		return msgType, nil, nil, fmt.Errorf("message type %v dose not exist", msgType)
	}

	buf.MoveReadOffset(LayerHdrSize)
	buf.MoveWriteOffset(msgLen)

	msg := regInfo.NewFunc()
	if err := msg.UnmarshalBinary(buf); err != nil {
		return msgType, nil, nil, fmt.Errorf("message type %v has failed to call UnmarshalBinary function, %v", msgType, err)
	}

	return msgType, msg, regInfo, nil
}

func LayerPacketMarshalBinary(msg message.IMessage, buf *nbuffer.BufferObject) error {
	buf.WriteBytes(byte(msg.GetType()), 0, 0)
	if err := msg.MarshalBinary(buf); err != nil {
		return err
	}

	// update message length
	msgLength := uint16(buf.GetWriteLength() - LayerHdrSize)
	binary.BigEndian.PutUint16(buf.GetRangeBytes(1, 2), msgLength)

	return nil
}
