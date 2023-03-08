package message

import (
	"github.com/whaoinfo/go-box/mapping"
	"github.com/whaoinfo/go-box/nbuffer"
	"github.com/whaoinfo/macro-UDP/message/element"
	"github.com/whaoinfo/macro-UDP/message/wrap"
)

func NewUploadInterfacePacketMessage() *UploadInterfacePacketMessage {
	return &UploadInterfacePacketMessage{}
}

type UploadInterfacePacketMessage struct {
	PodNameElement       *element.PodNameElement
	ServiceNameElement   *element.ServiceNameElement
	InterfaceNameElement *element.InterfaceNameElement
	PayloadPacketElement *element.PayloadPacketElement
}

func (t *UploadInterfacePacketMessage) MarshalBinary(buf *nbuffer.BufferObject) error {
	callArgs := mapping.ListToValues(buf)
	return mapping.ScanAllFields(t, wrap.CallFieldMarshalBinary, callArgs)
}

func (t *UploadInterfacePacketMessage) UnmarshalBinary(buf *nbuffer.BufferObject) error {
	callArgs := mapping.ListToValues(buf)
	return mapping.ScanAllFields(t, wrap.CallFieldUnmarshalBinary, callArgs)
}
