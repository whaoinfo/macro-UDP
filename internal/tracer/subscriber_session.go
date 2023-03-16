package tracer

import (
	"bytes"
	"fmt"
	"github.com/whaoinfo/go-box/logger"
	"github.com/whaoinfo/macro-UDP/internal/configmodel"
	message "github.com/whaoinfo/macro-UDP/internal/message"
)

type SubscriberFilter struct {
	IMSI  []byte
	MSIDN []byte
	IMEI  []byte
}

func (t *SubscriberFilter) Check(iMsg message.IMessage) bool {
	msg := iMsg.(*message.UploadSubscriberPacketMessage)
	if len(t.IMEI) > 0 {
		if msg.SubscriberIdentitiesElement.IMSIElement != nil &&
			bytes.Compare(msg.SubscriberIdentitiesElement.IMSIElement.Value, t.IMSI) != 0 {
			return false
		}
	}

	if len(t.MSIDN) > 0 {
		if msg.SubscriberIdentitiesElement.MSISDNElement != nil &&
			bytes.Compare(msg.SubscriberIdentitiesElement.MSISDNElement.Value, t.MSIDN) != 0 {
			return false
		}
	}

	if len(t.IMEI) > 0 {
		if msg.SubscriberIdentitiesElement.IMEIElement != nil &&
			bytes.Compare(msg.SubscriberIdentitiesElement.IMEIElement.Value, t.IMEI) != 0 {
			return false
		}
	}

	return true
}

func NewSubscriberFilters(cfg *configmodel.ConfigTraceSessionModel) ([]IFilter, error) {
	var retList []IFilter
	for n, elem := range cfg.SubscriberList {
		if elem.IMSI == "" && elem.MSIDN == "" && elem.IMEI == "" {
			return nil, fmt.Errorf("the subscriber list element of index %d is empty", n)
		}

		filter := &SubscriberFilter{}
		if elem.IMSI != "" {
			filter.IMSI = []byte(elem.IMSI)
		}
		if elem.MSIDN != "" {
			filter.MSIDN = []byte(elem.MSIDN)
		}
		if elem.IMEI != "" {
			filter.IMEI = []byte(elem.IMEI)
		}

		logger.DebugFmt("Parse subscriber filter config, TraceSessionId: %v, Disable: %v, EndTime: %v, IMEI: %v, MSIDN: %v, "+
			"IMEI: %v",
			cfg.TraceSessionId, cfg.Disable, cfg.EndTime, elem.IMSI, elem.MSIDN, elem.IMEI)
		retList = append(retList, filter)
	}
	return retList, nil
}

//func (t *SubscriberSession) match(iMsg message.IMessage) bool {
//	if t.disable {
//		return false
//	}
//	if t.endTimestamp >= ctime.CurrentTimestamp() {
//		return false
//	}
//
//	msg := iMsg.(*message.UploadSubscriberPacketMessage)
//	if msg.SubscriberIdentitiesElement == nil {
//		return false
//	}
//
//	for _, filter := range t.filters {
//		if len(filter.IMEI) > 0 {
//			if msg.SubscriberIdentitiesElement.IMSIElement != nil &&
//				bytes.Compare(msg.SubscriberIdentitiesElement.IMSIElement.Value, filter.IMSI) != 0 {
//				continue
//			}
//		}
//
//		if len(filter.MSIDN) > 0 {
//			if msg.SubscriberIdentitiesElement.MSISDNElement != nil &&
//				bytes.Compare(msg.SubscriberIdentitiesElement.MSISDNElement.Value, filter.MSIDN) != 0 {
//				continue
//			}
//		}
//
//		if len(filter.IMEI) > 0 {
//			if msg.SubscriberIdentitiesElement.IMEIElement != nil &&
//				bytes.Compare(msg.SubscriberIdentitiesElement.IMEIElement.Value, filter.IMEI) != 0 {
//				continue
//			}
//		}
//		return true
//	}
//
//	return false
//}
