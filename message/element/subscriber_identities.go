package element

import (
	"encoding/binary"
	"github.com/whaoinfo/go-box/mapping"
	"github.com/whaoinfo/go-box/nbuffer"
	"github.com/whaoinfo/macro-UDP/message/wrap"
	"reflect"
)

type SubscriberIdentitiesElement struct {
	IMSIElement   *IMSIElement
	IMEIElement   *IMEIElement
	MSISDNElement *MSISDNElement
	UEIPV4Element *UEIPV4Element
	UEIPV6Element *UEIPV6Element
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
