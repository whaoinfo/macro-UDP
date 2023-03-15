package tracer

import (
	"bytes"
	"errors"
	"github.com/whaoinfo/go-box/ctime"
	"github.com/whaoinfo/go-box/logger"
	"github.com/whaoinfo/macro-UDP/internal/configmodel"
	"github.com/whaoinfo/macro-UDP/internal/message"
	"strings"
)

var (
	serviceNameAllFilter = []byte("all")
)

type FilterTraceInterface struct {
	PodName     []byte
	ServiceName []byte
	Interfaces  [][]byte
}

func NewInterfaceSession() ISession {
	return &InterfaceSession{}
}

type InterfaceSession struct {
	id           SessionID
	disable      bool
	endTimestamp int64
	filters      []*FilterTraceInterface
}

func (t *InterfaceSession) initialize(cfg *configmodel.ConfigTraceSessionModel) error {
	if cfg.TraceSessionId == "" {
		return errors.New("the traceSessionId field in the configuration is empty")
	}
	if cfg.EndTime == "" {
		return errors.New("the endTime field in the configuration is empty")
	}

	t.id = SessionID(cfg.TraceSessionId)
	t.disable = cfg.Disable
	//t.endTime = cfg.EndTime
	for _, elem := range cfg.InterfaceList {
		if elem.PodName == "" && elem.ServiceName == "" {
			return errors.New("the podName and serviceName fields in the configuration are empty")
		}
		if len(elem.Interfaces) <= 0 {
			return errors.New("the interfaces field in the configuration is empty")
		}

		interfaceFilter := &FilterTraceInterface{
			PodName:     []byte(elem.PodName),
			ServiceName: []byte(elem.ServiceName),
		}

		for _, netIf := range elem.Interfaces {
			interfaceFilter.Interfaces = append(interfaceFilter.Interfaces, []byte(netIf))
		}

		logger.DebugFmt("Parse config, TraceSessionId: %v, Disable: %v, EndTime: %v, PodName: %v, ServiceName: %v, "+
			"Interfaces: %v",
			cfg.TraceSessionId, cfg.Disable, cfg.EndTime, elem.PodName, elem.ServiceName, strings.Join(elem.Interfaces, ","))
		t.filters = append(t.filters, interfaceFilter)
	}

	return nil
}

func (t *InterfaceSession) getID() SessionID {
	return t.id
}

func (t *InterfaceSession) match(iMsg message.IMessage) bool {
	if t.disable {
		return false
	}
	if t.endTimestamp >= ctime.CurrentTimestamp() {
		return false
	}

	msg := iMsg.(*message.UploadInterfacePacketMessage)
	if msg.InterfaceNameElement == nil || len(msg.InterfaceNameElement.Value) <= 0 {
		return false
	}

	filters := t.filters
	for _, filter := range filters {
		if len(filter.ServiceName) > 0 {
			if msg.ServiceNameElement == nil || len(msg.ServiceNameElement.Value) <= 0 {
				continue
			}

			if bytes.Compare(filter.ServiceName, msg.ServiceNameElement.Value) != 0 {
				continue
			}
		}

		if len(filter.PodName) > 0 && bytes.Compare(filter.PodName, serviceNameAllFilter) != 0 {
			if msg.PodNameElement == nil || len(msg.PodNameElement.Value) <= 0 {
				continue
			}

			if bytes.Compare(filter.PodName, msg.PodNameElement.Value) != 0 {
				continue
			}
		}

		for _, d := range filter.Interfaces {
			if bytes.Compare(d, msg.InterfaceNameElement.Value) == 0 {
				return true
			}
		}
	}

	return false
}
