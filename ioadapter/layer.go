package ioadapter

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/whaoinfo/go-box/logger"
	"github.com/whaoinfo/go-box/nbuffer"
)

const (
	MsgTypeSize   = 1
	MsgLengthSize = 2
	LayerHdrSize  = MsgTypeSize + MsgLengthSize
)

func NewLayerPacket() *LayerPacket {
	return &LayerPacket{}
}

type LayerPacketHeader struct {
	MsgType   MsgType
	MsgLength uint16
}

type LayerPacket struct {
	LayerPacketHeader
	Msg IMessage

	buf *nbuffer.BufferObject
}

func (t *LayerPacket) UnmarshalBinary(buf *nbuffer.BufferObject) error {
	bufBytes := buf.Bytes()
	layerHdrBuf := bufBytes[0:LayerHdrSize]
	// parse message type
	msgType := MsgType(layerHdrBuf[0])
	logger.DebugFmt("Process message, type: %v", msgType)

	// parse message length
	msgLen := int(binary.BigEndian.Uint16(layerHdrBuf[1:]))
	// Check length
	if (buf.GetWriteLength() - LayerHdrSize) < msgLen {
		return errors.New("message length does not match")
	}

	t.MsgType = msgType
	t.MsgLength = uint16(msgLen)

	// Find message type
	newMsgFunc, exist := newMessageFuncMap[msgType]
	if !exist {
		return fmt.Errorf("message type %v is not registered", msgType)
	}

	buf.MoveReadOffset(LayerHdrSize)
	buf.MoveWriteOffset(msgLen)

	msg := newMsgFunc()
	if err := msg.UnmarshalBinary(buf); err != nil {
		return fmt.Errorf("the %v type message has failed to call UnmarshalBinary function, %v", msgType, err)
	}

	t.buf = buf
	return nil
}

func (t *LayerPacket) MarshalBinary(buf *nbuffer.BufferObject) error {
	buf.WriteBytes(byte(t.MsgType), 0, 0)

	if err := t.Msg.MarshalBinary(buf); err != nil {
		return err
	}

	// update message length
	t.MsgLength = uint16(buf.GetWriteLength() - LayerHdrSize)
	binary.BigEndian.PutUint16(buf.GetRangeBytes(1, 2), t.MsgLength)

	return nil
}
