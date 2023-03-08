package message

import (
	"github.com/whaoinfo/go-box/mapping"
	"github.com/whaoinfo/go-box/nbuffer"
	"github.com/whaoinfo/macro-UDP/message/element"
	"github.com/whaoinfo/macro-UDP/message/wrap"
)

func NewUploadSubscriberPacketMessage() *UploadSubscriberPacketMessage {
	return &UploadSubscriberPacketMessage{}
}

type UploadSubscriberPacketMessage struct {
	NFInstanceIDElement         *element.NFInstanceIDElement
	SubscriberIdentitiesElement *element.SubscriberIdentitiesElement
	ExtensionFieldElement       *element.ExtensionFieldElement
	PayloadPacketElement        *element.PayloadPacketElement
}

func (t *UploadSubscriberPacketMessage) MarshalBinary(buf *nbuffer.BufferObject) error {
	callArgs := mapping.ListToValues(buf)
	markFilter := false
	return mapping.ScanAllFields(t, wrap.CallFieldMarshalBinary, callArgs, &markFilter)
}

func (t *UploadSubscriberPacketMessage) UnmarshalBinary(buf *nbuffer.BufferObject) error {
	callArgs := mapping.ListToValues(buf)
	return mapping.ScanAllFields(t, wrap.CallFieldUnmarshalBinary, callArgs)
}
