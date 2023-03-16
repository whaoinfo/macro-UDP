package tracer

import (
	"github.com/whaoinfo/go-box/logger"
	cfgmd "github.com/whaoinfo/macro-UDP/internal/configmodel"
	"github.com/whaoinfo/macro-UDP/internal/message"
	frame "github.com/whaoinfo/macro-UDP/pkg/gicframe"
)

type ComponentKW struct {
	MaxQueueNum    int `json:"max_queue_num"`
	MaxQueueLength int `json:"max_queue_length"`
}

type Component struct {
	frame.BaseComponent
	tracerMap map[tracerType]ITracer
	//kw *ComponentKW
}

func (t *Component) Initialize(kw frame.IComponentKW) error {
	t.tracerMap = make(map[tracerType]ITracer)
	//t.kw = kw.(*ComponentKW)
	kwArgs := kw.(*ComponentKW)
	cfg := cfgmd.GetTracingConfig()

	for tracerTpy, f := range regNewTracerFuncMap {
		tracer := f()
		if !t.bindMessages(tracer) {
			logger.WarnFmt("Tracer type %v has failed to bind message, %v", tracerTpy)
			continue
		}

		if err := tracer.setQueueInfo(kwArgs.MaxQueueNum, kwArgs.MaxQueueLength); err != nil {
			logger.WarnFmt("Tracer type %v has failed to set queue info, %v", tracerTpy, err)
			continue
		}

		t.addTracerSessions(tracer, cfg)
		t.tracerMap[tracerTpy] = tracer
		logger.InfoFmt("Added tracer type %v to the component, MaxQueueNum: %d, MaxQueueLength: %d",
			tracerTpy, kwArgs.MaxQueueNum, kwArgs.MaxQueueLength)
	}

	return nil
}

func (t *Component) Start() error {
	for _, tracer := range t.tracerMap {
		if err := tracer.start(); err != nil {
			logger.WarnFmt("Tracer type %v has failed to start, %v", err)
		}
	}
	return nil
}

func (t *Component) Stop() error {
	for _, tracer := range t.tracerMap {
		if err := tracer.stop(); err != nil {
			logger.WarnFmt("Tracer type %v has failed to stop, %v", err)
		}
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
		msgInfo.Handle = t.traceMessageHandle
		msgInfo.Args = []interface{}{tracer.getType()}
		logger.InfoFmt("Tracer type %v is already bound to message type %v", tracer.getType(), msgTpy)
	}

	return true
}

func (t *Component) addTracerSessions(tracer ITracer, cfg *cfgmd.ConfigTracingModel) {
	for _, sessCfg := range cfg.TraceSessionList {
		if !tracer.checkConfigType(&sessCfg) {
			continue
		}
		if err := tracer.addSessionByConfig(&sessCfg); err != nil {
			logger.WarnFmt("Tracer type %v has failed to add the %v session, %v",
				tracer.getType(), sessCfg.TraceSessionId, err)
		}
	}
}

func (t *Component) traceMessageHandle(ctx *message.HandleContext, args ...interface{}) {
	tpy := args[0].(tracerType)
	handleOk := false
	defer func() {
		if !handleOk {
			// todo: recycle ctx
		}
	}()

	logger.AllFmt("Tracing message, TracerType: %v", tpy)
	if ctx == nil {
		logger.WarnFmt("Failed to trace on type %v, the ctx is a nil pointer", tpy)
		return
	}
	if ctx.Msg == nil {
		logger.WarnFmt("Failed to trace on type %v, the ctx.msg is a nil pointer", tpy)
		return
	}
	if ctx.Buf == nil {
		logger.WarnFmt("Failed to trace on type %v, the ctx.buf is a nil pointer", tpy)
		return
	}

	// put messages
	tracer := t.tracerMap[tpy]
	if tracer == nil {
		logger.WarnFmt("Failed to trace on type %v, the type dose not exist", tpy)
		return
	}

	if !tracer.traceMessage(ctx) {
		logger.WarnFmt("Failed to trace on type %v, the type dose not exist", tpy)
		return
	}

	handleOk = true
	logger.AllFmt("traced message, TracerType: %v, MsgType: %v", tracer.getType(), ctx.Msg.GetType())

}
