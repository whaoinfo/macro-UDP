package tracer

import (
	"github.com/whaoinfo/macro-UDP/internal/configmodel"
	"github.com/whaoinfo/macro-UDP/internal/message"
)

func NewInterfaceSession() ISession {
	return &InterfaceSession{}
}

type InterfaceSession struct {
	id      SessionID
	disable bool

	endTimestamp int64
}

func (t *InterfaceSession) initialize(cfg *configmodel.ConfigTraceSessionModel) error {
	t.id = SessionID(cfg.TraceSessionId)
	t.disable = cfg.Disable
	//t.endTime = cfg.EndTime

	return nil
}

func (t *InterfaceSession) getID() SessionID {
	return t.id
}

func (t *InterfaceSession) match(msg message.IMessage) bool {
	return false
}
