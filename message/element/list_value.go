package element

import (
	"fmt"
	"github.com/whaoinfo/go-box/mapping"
	"github.com/whaoinfo/go-box/nbuffer"
)

func ListValueUnmarshalBinary(lenSize int, buf *nbuffer.BufferObject, value *[]byte) error {
	bytesToNumFunc := mapping.GetBytesToNumberFunc(lenSize)
	if bytesToNumFunc == nil {
		return fmt.Errorf("the %v length size type is invalid", lenSize)
	}

	lenBytes := buf.Read(lenSize)
	length := bytesToNumFunc(lenBytes)
	*value = buf.Read(length)
	return nil
}

func ListValueMarshalBinary(lenSize int, buf *nbuffer.BufferObject, value []byte) error {
	length := len(value)
	if length <= 0 {
		return nil
	}

	numToBytesFunc := mapping.GetNumberToBytesFunc(lenSize)
	if numToBytesFunc == nil {
		return fmt.Errorf("the %v length size type is invalid", lenSize)
	}

	nextWriteBytes := buf.GetNextWriteBytes()
	numToBytesFunc(length, nextWriteBytes)
	buf.MoveWriteOffset(lenSize)
	return buf.Write(value)
}

type ListValueU8Element struct {
	Value []byte
}

func (t *ListValueU8Element) UnmarshalBinary(buf *nbuffer.BufferObject) error {
	return ListValueUnmarshalBinary(mapping.UINT8Size, buf, &t.Value)
}

func (t *ListValueU8Element) MarshalBinary(buf *nbuffer.BufferObject) error {
	return ListValueMarshalBinary(mapping.UINT8Size, buf, t.Value)
}

type ListValueU16Element struct {
	Value []byte
}

func (t *ListValueU16Element) UnmarshalBinary(buf *nbuffer.BufferObject) error {
	return ListValueUnmarshalBinary(mapping.UINT16Size, buf, &t.Value)
}

func (t *ListValueU16Element) MarshalBinary(buf *nbuffer.BufferObject) error {
	return ListValueMarshalBinary(mapping.UINT16Size, buf, t.Value)
}
