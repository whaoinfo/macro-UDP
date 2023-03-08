package ioadapter

import "github.com/whaoinfo/go-box/nbuffer"

type MsgType uint8

type NewMessageFunc func() IMessage

type IMessage interface {
	MarshalBinary(buf *nbuffer.BufferObject) error
	UnmarshalBinary(buf *nbuffer.BufferObject) error
}

var (
	newMessageFuncMap = make(map[MsgType]NewMessageFunc)
)

func RegisterMessage(msgType MsgType, newFunc NewMessageFunc) {
	newMessageFuncMap[msgType] = newFunc
}
