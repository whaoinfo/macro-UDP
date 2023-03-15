package tracer

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/whaoinfo/go-box/ctime"
	"github.com/whaoinfo/go-box/logger"
	"github.com/whaoinfo/macro-UDP/internal/configmodel"
	message "github.com/whaoinfo/macro-UDP/internal/message"
)

type FilterTraceSubscriber struct {
	IMSI  []byte
	MSIDN []byte
	IMEI  []byte
}

func NewSubscriberSession() ISession {
	return &SubscriberSession{}
}

type SubscriberSession struct {
	id           SessionID
	disable      bool
	endTimestamp int64
	filters      []*FilterTraceSubscriber
}

func (t *SubscriberSession) initialize(cfg *configmodel.ConfigTraceSessionModel) error {
	if cfg.TraceSessionId == "" {
		return errors.New("the traceSessionId field in the configuration is empty")
	}
	if cfg.EndTime == "" {
		return errors.New("the endTime field in the configuration is empty")
	}

	t.id = SessionID(cfg.TraceSessionId)
	t.disable = cfg.Disable

	for n, elem := range cfg.SubscriberList {
		if elem.IMSI == "" && elem.MSIDN == "" && elem.IMEI == "" {
			return fmt.Errorf("the subscriber list element of index %d is empty", n)
		}

		filter := &FilterTraceSubscriber{}
		if elem.IMSI != "" {
			filter.IMSI = []byte(elem.IMSI)
		}
		if elem.MSIDN != "" {
			filter.MSIDN = []byte(elem.MSIDN)
		}
		if elem.IMEI != "" {
			filter.IMEI = []byte(elem.IMEI)
		}

		logger.DebugFmt("Parse config, TraceSessionId: %v, Disable: %v, EndTime: %v, IMEI: %v, MSIDN: %v, "+
			"IMEI: %v",
			cfg.TraceSessionId, cfg.Disable, cfg.EndTime, elem.IMSI, elem.MSIDN, elem.IMEI)
		t.filters = append(t.filters, filter)
	}

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

	msg := iMsg.(*message.UploadSubscriberPacketMessage)
	if msg.SubscriberIdentitiesElement == nil {
		return false
	}

	for _, filter := range t.filters {
		if len(filter.IMEI) > 0 {
			if msg.SubscriberIdentitiesElement.IMSIElement != nil &&
				bytes.Compare(msg.SubscriberIdentitiesElement.IMSIElement.Value, filter.IMSI) != 0 {
				continue
			}
		}

		if len(filter.MSIDN) > 0 {
			if msg.SubscriberIdentitiesElement.MSISDNElement != nil &&
				bytes.Compare(msg.SubscriberIdentitiesElement.MSISDNElement.Value, filter.MSIDN) != 0 {
				continue
			}
		}

		if len(filter.IMEI) > 0 {
			if msg.SubscriberIdentitiesElement.IMEIElement != nil &&
				bytes.Compare(msg.SubscriberIdentitiesElement.IMEIElement.Value, filter.IMEI) != 0 {
				continue
			}
		}

		return true
	}

	return false
}
