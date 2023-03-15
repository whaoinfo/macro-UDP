package tracer

import (
	"github.com/whaoinfo/go-box/logger"
	"github.com/whaoinfo/go-box/nbuffer"
	configmodel "github.com/whaoinfo/macro-UDP/internal/configmodel"
	"github.com/whaoinfo/macro-UDP/internal/message"
	frame "github.com/whaoinfo/macro-UDP/pkg/gicframe"
)

type Component struct {
	frame.BaseComponent
	tracerMap map[tracerType]ITracer
}

func (t *Component) Initialize(kw frame.IComponentKW) error {
	t.tracerMap = make(map[tracerType]ITracer)
	cfg := configmodel.GetConfigModel().Tracer

	for tracerTpy, f := range regNewTracerFuncMap {
		tracer := f()
		if !t.bindMessages(tracer) {
			logger.WarnFmt("Tracer type %v has failed to bind message, %v", tracerTpy)
			continue
		}

		t.addTracerSessions(tracer, cfg)

		t.tracerMap[tracerTpy] = tracer
		logger.InfoFmt("Added tracer type %v to the component", tracerTpy)
	}

	return nil
}

func (t *Component) bindMessages(tracer ITracer) bool {
	for _, msgTpy := range tracer.getRefMessageTypes() {
		msgInfo := message.GetRegisterMessageInfo(msgTpy)
		if msgInfo == nil {
			logger.WarnFmt("Tracer type %v gets message type %v does not exist", tracer.getType(), msgTpy)
			return false
		}

		if msgInfo.Handle != nil {
			logger.WarnFmt("Message type %v is repeatedly bound by tracer type %v", msgTpy, tracer.getType())
			return false
		}
		msgInfo.Handle = t.traceMessage
		msgInfo.Args = []interface{}{tracer.getType()}
		logger.InfoFmt("Tracer type %v is already bound to message type %v", tracer.getType(), msgTpy)
	}

	return true
}

func (t *Component) addTracerSessions(tracer ITracer, cfg configmodel.ConfigTracerModel) {
	for _, sessCfg := range cfg {
		if !tracer.checkConfigType(&sessCfg) {
			continue
		}
		if err := tracer.addSessionByConfig(&sessCfg); err != nil {
			logger.WarnFmt("Tracer type %v has failed to add the %v session, %v",
				tracer.getType(), sessCfg.TraceSessionId, err)
		}
	}

}

func (t *Component) traceMessage(msg message.IMessage, buf *nbuffer.BufferObject, args ...interface{}) {
	tpy := args[0].(tracerType)
	if msg == nil {
		logger.WarnFmt("Failed to trace message, the message is a nil pointer. tracer type: %v", tpy)
	}

	// match sessions
	tracer := t.tracerMap[tpy]
	if tracer == nil {
		logger.WarnFmt("Failed to trace message type %v, tracer type %v does not exist", msg.GetType(), tpy)
		return
	}
	// put messages to storage
	for _, sess := range tracer.matchSessions(msg) {
		sess = sess
	}

	// Put msg to storage

}
