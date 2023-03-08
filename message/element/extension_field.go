package element

import (
	"github.com/whaoinfo/go-box/mapping"
	"github.com/whaoinfo/go-box/nbuffer"
)

type ExtensionFieldElement struct {
	Flag             uint8
	ProtocolType     uint8
	MicroserviceType uint8
}

func (t *ExtensionFieldElement) UnmarshalBinary(bufObj *nbuffer.BufferObject) error {
	t.Flag = bufObj.Read(mapping.UINT8Size)[0]
	if t.Flag <= 0 {
		return nil
	}

	for n, elem := range []*uint8{&t.ProtocolType, &t.MicroserviceType} {
		if (t.Flag & (1 << n)) != 0 {
			*elem = bufObj.Read(mapping.UINT8Size)[0]
		}
	}

	return nil
}

func (t *ExtensionFieldElement) MarshalBinary(bufObj *nbuffer.BufferObject) error {
	bufObj.WriteBytes(t.Flag)
	if t.Flag <= 0 {
		return nil
	}

	for n, elem := range []uint8{t.ProtocolType, t.MicroserviceType} {
		if (t.Flag & (1 << n)) != 0 {
			bufObj.WriteBytes(elem)
		}
	}

	return nil
}
