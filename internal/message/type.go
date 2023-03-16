package message

import (
	"github.com/whaoinfo/go-box/nbuffer"
)

const (
	UploadSubscriberPacketMessageType MsgType = iota + 1
	UploadInterfacePacketMessageType
)

type MsgType uint8
type NewFunc func() IMessage
type Handle func(ctx *HandleContext, args ...interface{})

type HandleContext struct {
	Buf *nbuffer.BufferObject
	Msg IMessage
}

type IMessage interface {
	MarshalBinary(buf *nbuffer.BufferObject) error
	UnmarshalBinary(buf *nbuffer.BufferObject) error
	GetType() MsgType
}

type RegMessageInfo struct {
	NewFunc NewFunc
	Handle  Handle
	Args    []interface{}
}

var (
	regMessageInfo = map[MsgType]*RegMessageInfo{
		UploadSubscriberPacketMessageType: {
			NewFunc: func() IMessage {
				return &UploadSubscriberPacketMessage{}
			},
		},
		UploadInterfacePacketMessageType: {
			NewFunc: func() IMessage {
				return &UploadInterfacePacketMessage{}
			},
		},
	}
)

func GetRegisterMessageInfo(tpy MsgType) *RegMessageInfo {
	return regMessageInfo[tpy]
}
