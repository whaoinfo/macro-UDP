package bufferelement

import (
	"errors"
	"github.com/whaoinfo/go-box/nbuffer"
	"net"
)

type UEIPV4Element struct {
	IP net.IP
}

func (t *UEIPV4Element) UnmarshalBinary(buf *nbuffer.BufferObject) error {
	d := buf.Read(net.IPv4len)
	nip := net.IP(d).To4()
	if nip == nil {
		return errors.New("invalid IP V4 format")
	}

	t.IP = nip
	return nil
}

func (t *UEIPV4Element) MarshalBinary(buf *nbuffer.BufferObject) error {
	if t.IP == nil {
		return errors.New("invalid IP V4 format")
	}

	d := buf.GetNextWriteBytes()
	copy(d, t.IP.To4())
	buf.MoveWriteOffset(net.IPv4len)
	return nil
}

type UEIPV6Element struct {
	IP net.IP
}

func (t *UEIPV6Element) UnmarshalBinary(buf *nbuffer.BufferObject) error {
	d := buf.Read(net.IPv6len)
	nip := net.IP(d).To16()
	if nip == nil {
		return errors.New("invalid IP V6 format")
	}

	t.IP = nip
	return nil
}

func (t *UEIPV6Element) MarshalBinary(buf *nbuffer.BufferObject) error {
	if t.IP == nil {
		return errors.New("invalid IP V6 format")
	}

	d := buf.GetNextWriteBytes()
	copy(d, t.IP.To16())
	buf.MoveWriteOffset(net.IPv6len)
	return nil
}
