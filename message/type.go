package message

import "github.com/whaoinfo/macro-UDP/ioadapter"

const (
	UploadSubscriberPacketMessageType = iota + 1
	UploadInterfacePacketMessageType
)

func init() {
	ioadapter.RegisterMessage(UploadSubscriberPacketMessageType, func() ioadapter.IMessage {
		return NewUploadSubscriberPacketMessage()
	})
	ioadapter.RegisterMessage(UploadInterfacePacketMessageType, func() ioadapter.IMessage {
		return NewUploadInterfacePacketMessage()
	})
}
