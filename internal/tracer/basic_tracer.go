package tracer

import (
	"errors"
	"fmt"
	"github.com/whaoinfo/go-box/logger"
	"github.com/whaoinfo/macro-UDP/internal/configmodel"
	"github.com/whaoinfo/macro-UDP/internal/define"
	"github.com/whaoinfo/macro-UDP/internal/message"
)

func NewBasicTracer(tpy tracerType, newSession NewSessionFunc, checkConfigTypeFunc checkConfigTypeFunc,
	refMsgTypes []message.MsgType) *BasicTracer {

	return &BasicTracer{
		tpy:          tpy,
		newSession:   newSession,
		refMsgTypes:  refMsgTypes,
		sessionGroup: &sessionGroup{elems: make(map[define.SessionID]ISession)},
		//sessionMap:          make(map[SessionID]ISession),
		checkConfigTypeFunc: checkConfigTypeFunc,
	}
}

type checkConfigTypeFunc func(cfg *configmodel.ConfigTraceSessionModel) bool

type sessionGroup struct {
	elems map[define.SessionID]ISession
}

type BasicTracer struct {
	newSession   NewSessionFunc
	sessionGroup *sessionGroup
	//lock sync.RWMutex
	refMsgTypes         []message.MsgType
	tpy                 tracerType
	checkConfigTypeFunc checkConfigTypeFunc
	maxQueueNum         int
	maxQueueLength      int
}

func (t *BasicTracer) setQueueInfo(maxNum, maxlength int) error {
	if maxNum <= 0 || maxlength <= 0 {
		return errors.New("maxNum <= 0 || maxlength <= 0")
	}

	t.maxQueueNum = maxNum
	t.maxQueueLength = maxlength
	return nil
}

func (t *BasicTracer) getType() tracerType {
	return t.tpy
}

func (t *BasicTracer) checkConfigType(cfg *configmodel.ConfigTraceSessionModel) bool {
	if t.checkConfigTypeFunc == nil {
		logger.WarnFmt("The CheckConfigType function of tracer type %v is not implemented", t.tpy)
		return false
	}

	return t.checkConfigTypeFunc(cfg)
}

func (t *BasicTracer) getRefMessageTypes() []message.MsgType {
	return t.refMsgTypes
}

func (t *BasicTracer) start() error {
	for _, sess := range t.sessionGroup.elems {
		if err := sess.start(); err != nil {
			logger.WarnFmt("%v session has failed to start, %v", sess.getID(), err)
			continue
		}
		logger.DebugFmt("%v session has started", sess.getID())
	}
	return nil
}

func (t *BasicTracer) stop() error {
	for _, sess := range t.sessionGroup.elems {
		if err := sess.stop(); err != nil {
			logger.WarnFmt("%v session has failed to stop, %v", err)
			continue
		}
		logger.DebugFmt("%v session has stopped", sess.getID())
	}

	return nil
}

func (t *BasicTracer) addSessionByConfig(cfg *configmodel.ConfigTraceSessionModel) error {
	sess := t.newSession()
	if err := sess.initialize(t.maxQueueNum, t.maxQueueLength, cfg); err != nil {
		return fmt.Errorf("failed to initialize, %v", err)
	}

	t.sessionGroup.elems[sess.getID()] = sess
	logger.InfoFmt("Added %v session of tracer type %v", sess.getID(), t.getType())
	return nil
}

func (t *BasicTracer) traceMessage(ctx *message.HandleContext) bool {
	sessGroup := t.sessionGroup
	if sessGroup == nil {
		logger.WarnFmt("The sessionGroup field of tracer type %v is a nil pointer")
		return false
	}

	putOk := false
	for _, sess := range sessGroup.elems {
		if !sess.match(ctx.Msg) {
			logger.AllFmt("Matched %v session, TracerType: %v, MessageType: %v",
				sess.getID(), t.getType(), ctx.Msg.GetType())
			continue
		}
		if sess.putMessageContext(ctx) {
			putOk = true
		}
	}

	return putOk
}
