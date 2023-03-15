package message

import (
	"github.com/whaoinfo/go-box/mapping"
	"github.com/whaoinfo/go-box/nbuffer"
	"github.com/whaoinfo/macro-UDP/internal/message/wrap"
	"github.com/whaoinfo/macro-UDP/pkg/bufferelement"
)

type UploadInterfacePacketMessage struct {
	PodNameElement       *bufferelement.U8BytesValueElement
	ServiceNameElement   *bufferelement.U8BytesValueElement
	InterfaceNameElement *bufferelement.U8BytesValueElement
	PayloadPacketElement *bufferelement.U16ListValueElement
}

func (t *UploadInterfacePacketMessage) GetType() MsgType {
	return UploadInterfacePacketMessageType
}

func (t *UploadInterfacePacketMessage) MarshalBinary(buf *nbuffer.BufferObject) error {
	callArgs := mapping.ListToValues(buf)
	return mapping.ScanAllFields(t, wrap.CallFieldMarshalBinary, callArgs)
}

func (t *UploadInterfacePacketMessage) UnmarshalBinary(buf *nbuffer.BufferObject) error {
	callArgs := mapping.ListToValues(buf)
	return mapping.ScanAllFields(t, wrap.CallFieldUnmarshalBinary, callArgs)
}

func InterfaceMessageTrace(msg IMessage, buf *nbuffer.BufferObject) {

}
