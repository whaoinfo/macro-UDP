package tracer

import (
	"github.com/whaoinfo/macro-UDP/internal/configmodel"
	"github.com/whaoinfo/macro-UDP/internal/message"
)

type tracerType uint8
type SessionType uint8
type SessionID string
type SessionGroup map[SessionID]ISession
type NewSessionFunc func() ISession

const (
	subscriberTracerType tracerType = iota + 1
	interfaceTracerType
)

type ITracer interface {
	getType() tracerType
	checkConfigType(cfg *configmodel.ConfigTraceSessionModel) bool
	addSessionByConfig(cfg *configmodel.ConfigTraceSessionModel) error
	getRefMessageTypes() []message.MsgType
	matchSessions(msg message.IMessage) []ISession
}

var (
	regNewTracerFuncMap = map[tracerType]func() ITracer{
		subscriberTracerType: func() ITracer {
			return NewBaseTracer(subscriberTracerType,
				NewSubscriberSession,
				func(cfg *configmodel.ConfigTraceSessionModel) bool {
					return len(cfg.SubscriberList) > 0
				},
				[]message.MsgType{message.UploadSubscriberPacketMessageType})
		},

		interfaceTracerType: func() ITracer {
			return NewBaseTracer(interfaceTracerType,
				NewInterfaceSession,
				func(cfg *configmodel.ConfigTraceSessionModel) bool {
					return len(cfg.InterfaceList) > 0
				},
				[]message.MsgType{message.UploadInterfacePacketMessageType})
		},
	}
)

type ISession interface {
	initialize(*configmodel.ConfigTraceSessionModel) error
	getID() SessionID
	match(msg message.IMessage) bool
}
