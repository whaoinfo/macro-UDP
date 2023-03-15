package main

import (
	"flag"
	"github.com/whaoinfo/go-box/logger"
	"github.com/whaoinfo/macro-UDP/internal/bufferpool"
	"github.com/whaoinfo/macro-UDP/internal/ioadapter"
	message2 "github.com/whaoinfo/macro-UDP/internal/message"
	element2 "github.com/whaoinfo/macro-UDP/internal/message/element"
	"github.com/whaoinfo/macro-UDP/pkg/bufferelement"
	"io"
	"net"
	"os"
	"time"
)

func newMessage() message2.IMessage {
	msg := &message2.UploadSubscriberPacketMessage{}
	nfInstID := &bufferelement.U8BytesValueElement{}
	nfInstID.Value = []byte("f2b43e11-762a-49b9-b5d4-4c3452ecb8f9")
	msg.NFInstanceIDElement = nfInstID
	logger.DebugFmt("NFInstanceIDElement value: %v, len: %d", string(nfInstID.Value), len(nfInstID.Value))

	subIdentities := &element2.SubscriberIdentitiesElement{}
	imsiElem := &bufferelement.U8BytesValueElement{}
	imsiElem.Value = []byte("1111111111")
	subIdentities.IMSIElement = imsiElem
	logger.DebugFmt("IMSIElement value: %v, len: %d", string(imsiElem.Value), len(imsiElem.Value))
	//subIdentities.Flag = subIdentities.Flag | 1

	subIdentities.IMEIElement = &bufferelement.U8BytesValueElement{}
	subIdentities.IMEIElement.Value = []byte("2222222222")
	logger.DebugFmt("IMEIElement value: %v, len: %d", string(subIdentities.IMEIElement.Value), len(subIdentities.IMEIElement.Value))
	//subIdentities.Flag = subIdentities.Flag | (1 << 2)

	subIdentities.UEIPV4Element = &bufferelement.UEIPV4Element{}
	subIdentities.UEIPV4Element.IP = net.ParseIP("172.10.0.15")
	//subIdentities.Flag = subIdentities.Flag | (1 << 4)
	msg.SubscriberIdentitiesElement = subIdentities

	extFieldElement := &element2.ExtensionFieldElement{
		ProtocolType: 2, MicroserviceType: 3,
	}
	extFieldElement.Flag = 0x3
	msg.ExtensionFieldElement = extFieldElement

	subPacket := &bufferelement.U16ListValueElement{}
	subPacket.Value = []byte("ccc-vvv-bbb-nnn-mmm")
	logger.DebugFmt("PayloadPacketElement value: %v, len: %d", string(subPacket.Value), len(subPacket.Value))
	msg.PayloadPacketElement = subPacket

	return msg
}

func newWriter(addr string) (io.Writer, error) {
	writer, err := net.Dial("udp", addr)
	if err != nil {
		return nil, err
	}

	return writer, nil
}

func main() {
	addr := flag.String("addr", "127.0.0.1:7777", "udp_address=127.0.0.1:8010")
	flag.Parse()

	msg := newMessage()

	writer, err := newWriter(*addr)
	if err != nil {
		logger.ErrorFmt("NewWriter Failed, %v", err)
		os.Exit(1)
	}

	buf := bufferpool.AllocateBufferObject()
	if err := ioadapter.LayerPacketMarshalBinary(msg, buf); err != nil {
		logger.ErrorFmt("LayerPacketMarshalBinary Failed, %v", err)
		os.Exit(1)
	}

	_, writeErr := writer.Write(buf.Bytes())
	if writeErr != nil {
		logger.ErrorFmt("Write Failed, %v", writeErr)
		os.Exit(1)
	}

	time.Sleep(time.Second * 2)
	os.Exit(0)
}
