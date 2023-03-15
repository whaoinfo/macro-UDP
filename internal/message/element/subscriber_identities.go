package element

import (
	"encoding/binary"
	"github.com/whaoinfo/go-box/logger"
	"github.com/whaoinfo/go-box/mapping"
	"github.com/whaoinfo/go-box/nbuffer"
	"github.com/whaoinfo/macro-UDP/internal/message/wrap"
	"github.com/whaoinfo/macro-UDP/pkg/bufferelement"
	"reflect"
)

type SubscriberIdentitiesElement struct {
	IMSIElement   *bufferelement.U8BytesValueElement
	IMEIElement   *bufferelement.U8BytesValueElement
	MSISDNElement *bufferelement.U8BytesValueElement
	UEIPV4Element *bufferelement.UEIPV4Element
	UEIPV6Element *bufferelement.UEIPV6Element
}

func (t *SubscriberIdentitiesElement) UnmarshalBinary(bufObj *nbuffer.BufferObject) error {
	d := bufObj.Read(mapping.UINT16Size)
	flag := binary.BigEndian.Uint16(d)
	if flag <= 0 {
		return nil
	}

	scanErr := mapping.ScanAllFields(t, func(ownRV reflect.Value, index int, args ...interface{}) error {
		marked := (flag & (1 << index)) != 0
		if !marked {
			return nil
		}

		return wrap.CallFieldUnmarshalBinary(ownRV, index, args...)
	}, mapping.ListToValues(bufObj))

	if t.IMSIElement != nil {
		logger.AllFmt("SubscriberIdentitiesElement.IMSIElement UnmarshalBinary, value: %v, len: %v",
			string(t.IMSIElement.Value), len(t.IMSIElement.Value))
	}
	if t.IMEIElement != nil {
		logger.AllFmt("SubscriberIdentitiesElement.IMEIElement UnmarshalBinary, value: %v, len: %v",
			string(t.IMEIElement.Value), len(t.IMEIElement.Value))
	}
	if t.MSISDNElement != nil {
		logger.AllFmt("SubscriberIdentitiesElement.MSISDNElement UnmarshalBinary, value: %v, len: %v",
			string(t.MSISDNElement.Value), len(t.MSISDNElement.Value))
	}
	if t.UEIPV4Element != nil {
		logger.AllFmt("SubscriberIdentitiesElement.UEIPV4Element UnmarshalBinary, value: %v",
			t.UEIPV4Element.IP.String())
	}
	if t.UEIPV6Element != nil {
		logger.AllFmt("SubscriberIdentitiesElement.UEIPV6Element UnmarshalBinary, value: %v",
			t.UEIPV6Element.IP.String())
	}

	return scanErr
}

func (t *SubscriberIdentitiesElement) MarshalBinary(bufObj *nbuffer.BufferObject) error {
	var flag uint16
	writeIdx := bufObj.GetWriteLength()
	bufObj.MoveWriteOffset(mapping.UINT16Size)
	markFilter := false
	scanErr := mapping.ScanAllFields(t, func(ownRV reflect.Value, index int, args ...interface{}) error {
		if err := wrap.CallFieldMarshalBinary(ownRV, index, args...); err != nil {
			return err
		}
		if markFilter {
			return nil
		}

		flag = flag | (1 << index)
		return nil
	}, mapping.ListToValues(bufObj), &markFilter)

	if scanErr != nil {
		return scanErr
	}

	binary.BigEndian.PutUint16(bufObj.GetRangeBytes(writeIdx, mapping.UINT16Size), flag)
	return nil
}
