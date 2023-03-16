package tracer

import (
	"github.com/whaoinfo/macro-UDP/internal/configmodel"
	"github.com/whaoinfo/macro-UDP/internal/define"
	"github.com/whaoinfo/macro-UDP/internal/message"
)

type tracerType uint8

type NewSessionFunc func() ISession
type NewFiltersFunc func(cfg *configmodel.ConfigTraceSessionModel) ([]IFilter, error)

const (
	subscriberTracerType tracerType = iota + 1
	interfaceTracerType
)

type ISession interface {
	initialize(maxQueueNum, queueLength int, cfg *configmodel.ConfigTraceSessionModel) error
	getID() define.SessionID
	start() error
	stop() error
	match(msg message.IMessage) bool
	putMessageContext(ctx *message.HandleContext) bool
}

type ITracer interface {
	getType() tracerType
	setQueueInfo(maxNum, maxlength int) error
	start() error
	stop() error
	checkConfigType(cfg *configmodel.ConfigTraceSessionModel) bool
	addSessionByConfig(cfg *configmodel.ConfigTraceSessionModel) error
	getRefMessageTypes() []message.MsgType
	traceMessage(ctx *message.HandleContext) bool
}

type IFilter interface {
	Check(msg message.IMessage) bool
}

var (
	regNewTracerFuncMap = map[tracerType]func() ITracer{
		subscriberTracerType: func() ITracer {
			return NewBasicTracer(subscriberTracerType,
				func() ISession {
					return &BaseSession{newFilters: NewSubscriberFilters}
				},
				func(cfg *configmodel.ConfigTraceSessionModel) bool {
					return len(cfg.SubscriberList) > 0
				},
				[]message.MsgType{message.UploadSubscriberPacketMessageType})
		},

		interfaceTracerType: func() ITracer {
			return NewBasicTracer(interfaceTracerType,
				func() ISession {
					return &BaseSession{newFilters: NewInterfaceFilters}
				},
				func(cfg *configmodel.ConfigTraceSessionModel) bool {
					return len(cfg.InterfaceList) > 0
				},
				[]message.MsgType{message.UploadInterfacePacketMessageType})
		},
	}
)
