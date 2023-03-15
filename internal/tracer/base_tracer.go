package tracer

import (
	"fmt"
	"github.com/whaoinfo/go-box/logger"
	"github.com/whaoinfo/macro-UDP/internal/configmodel"
	"github.com/whaoinfo/macro-UDP/internal/message"
)

func NewBaseTracer(tpy tracerType, newSession NewSessionFunc, checkConfigTypeFunc checkConfigTypeFunc,
	refMsgTypes []message.MsgType) *BaseTracer {

	return &BaseTracer{
		tpy:                 tpy,
		newSession:          newSession,
		refMsgTypes:         refMsgTypes,
		sessionMap:          make(map[SessionID]ISession),
		checkConfigTypeFunc: checkConfigTypeFunc,
	}
}

type checkConfigTypeFunc func(cfg *configmodel.ConfigTraceSessionModel) bool

type BaseTracer struct {
	newSession          NewSessionFunc
	sessionMap          map[SessionID]ISession
	refMsgTypes         []message.MsgType
	tpy                 tracerType
	checkConfigTypeFunc checkConfigTypeFunc
}

func (t *BaseTracer) getType() tracerType {
	return t.tpy
}

func (t *BaseTracer) checkConfigType(cfg *configmodel.ConfigTraceSessionModel) bool {
	if t.checkConfigTypeFunc == nil {
		logger.WarnFmt("The CheckConfigType function of tracer type %v is not implemented", t.tpy)
		return false
	}

	return t.checkConfigTypeFunc(cfg)
}

func (t *BaseTracer) getRefMessageTypes() []message.MsgType {
	return t.refMsgTypes
}

func (t *BaseTracer) addSessionByConfig(cfg *configmodel.ConfigTraceSessionModel) error {
	sess := t.newSession()
	if err := sess.initialize(cfg); err != nil {
		return fmt.Errorf("failed to initialize, %v", err)
	}

	t.sessionMap[sess.getID()] = sess
	logger.InfoFmt("Added the %v session of tracer type %v", sess.getID(), t.getType())
	return nil
}

func (t *BaseTracer) matchSessions(msg message.IMessage) []ISession {
	var retList []ISession
	for _, sess := range t.sessionMap {
		if sess.match(msg) {
			retList = append(retList, sess)
			logger.AllFmt("Matched the %v session, TracerType: %v, MessageType: %v",
				sess.getID(), t.getType(), msg.GetType())
		}
	}

	return retList
}
