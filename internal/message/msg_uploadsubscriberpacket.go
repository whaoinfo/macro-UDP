package message

import (
	"github.com/whaoinfo/go-box/mapping"
	"github.com/whaoinfo/go-box/nbuffer"
	element2 "github.com/whaoinfo/macro-UDP/internal/message/element"
	"github.com/whaoinfo/macro-UDP/internal/message/wrap"
	"github.com/whaoinfo/macro-UDP/pkg/bufferelement"
)

type UploadSubscriberPacketMessage struct {
	NFInstanceIDElement         *bufferelement.U8BytesValueElement
	SubscriberIdentitiesElement *element2.SubscriberIdentitiesElement
	ExtensionFieldElement       *element2.ExtensionFieldElement
	PayloadPacketElement        *bufferelement.U16ListValueElement
}

func (t *UploadSubscriberPacketMessage) GetType() MsgType {
	return UploadSubscriberPacketMessageType
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
