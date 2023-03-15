package tracer

import (
	"github.com/whaoinfo/go-box/ctime"
	"github.com/whaoinfo/macro-UDP/internal/configmodel"
	message "github.com/whaoinfo/macro-UDP/internal/message"
)

func NewSubscriberSession() ISession {
	return &SubscriberSession{}
}

type SubscriberSession struct {
	id      SessionID
	disable bool

	endTimestamp int64
	idents       []configmodel.ConfigTraceSubscriberModel
}

func (t *SubscriberSession) initialize(cfg *configmodel.ConfigTraceSessionModel) error {

	t.id = SessionID(cfg.TraceSessionId)
	t.disable = cfg.Disable
	//t.endTimestamp = cfg.EndTime

	return nil
}

func (t *SubscriberSession) getID() SessionID {
	return t.id
}

func (t *SubscriberSession) match(iMsg message.IMessage) bool {
	if t.disable {
		return false
	}
	if t.endTimestamp >= ctime.CurrentTimestamp() {
		return false
	}

	if iMsg.GetType() == message.UploadSubscriberPacketMessageType {
		msg := iMsg.(*message.UploadSubscriberPacketMessage)
		for _, elem := range t.idents {
			// todo: optimize to bytes -> string
			if elem.IMSI != "" && string(msg.SubscriberIdentitiesElement.IMSIElement.Value) == elem.IMSI {
				return true
			}
			if elem.MSIDN != "" && string(msg.SubscriberIdentitiesElement.MSISDNElement.Value) == elem.MSIDN {
				return true
			}
			if elem.IMEI != "" && string(msg.SubscriberIdentitiesElement.IMEIElement.Value) == elem.IMEI {
				return true
			}
		}
	}

	return false
}
