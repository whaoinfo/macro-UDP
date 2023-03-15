package bufferelement

import (
	"fmt"
	"github.com/whaoinfo/go-box/mapping"
	"github.com/whaoinfo/go-box/nbuffer"
)

func BytesValueUnmarshalBinary(lenSize int, buf *nbuffer.BufferObject, value *[]byte) error {
	bytesToNumFunc := mapping.GetBytesToNumberFunc(lenSize)
	if bytesToNumFunc == nil {
		return fmt.Errorf("the %v length size type is invalid", lenSize)
	}

	lenBytes := buf.Read(lenSize)
	length := bytesToNumFunc(lenBytes)
	*value = buf.Read(length)
	return nil
}

func BytesValueMarshalBinary(lenSize int, buf *nbuffer.BufferObject, value []byte) error {
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

type U8BytesValueElement struct {
	Value []byte
}

func (t *U8BytesValueElement) UnmarshalBinary(buf *nbuffer.BufferObject) error {
	return BytesValueUnmarshalBinary(mapping.UINT8Size, buf, &t.Value)
}

func (t *U8BytesValueElement) MarshalBinary(buf *nbuffer.BufferObject) error {
	return BytesValueMarshalBinary(mapping.UINT8Size, buf, t.Value)
}

type U16ListValueElement struct {
	Value []byte
}

func (t *U16ListValueElement) UnmarshalBinary(buf *nbuffer.BufferObject) error {
	return BytesValueUnmarshalBinary(mapping.UINT16Size, buf, &t.Value)
}

func (t *U16ListValueElement) MarshalBinary(buf *nbuffer.BufferObject) error {
	return BytesValueMarshalBinary(mapping.UINT16Size, buf, t.Value)
}

type U32ListValueElement struct {
	Value []byte
}

func (t *U32ListValueElement) UnmarshalBinary(buf *nbuffer.BufferObject) error {
	return BytesValueUnmarshalBinary(mapping.UINT32Size, buf, &t.Value)
}

func (t *U32ListValueElement) MarshalBinary(buf *nbuffer.BufferObject) error {
	return BytesValueMarshalBinary(mapping.UINT32Size, buf, t.Value)
}

type U64ListValueElement struct {
	Value []byte
}

func (t *U64ListValueElement) UnmarshalBinary(buf *nbuffer.BufferObject) error {
	return BytesValueUnmarshalBinary(mapping.UINT64Size, buf, &t.Value)
}

func (t *U64ListValueElement) MarshalBinary(buf *nbuffer.BufferObject) error {
	return BytesValueMarshalBinary(mapping.UINT64Size, buf, t.Value)
}
