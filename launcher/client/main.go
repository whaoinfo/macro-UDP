package main

import (
	"flag"
	"github.com/whaoinfo/go-box/logger"
	"github.com/whaoinfo/macro-UDP/bufferpool"
	"github.com/whaoinfo/macro-UDP/ioadapter"
	"github.com/whaoinfo/macro-UDP/message"
	"github.com/whaoinfo/macro-UDP/message/element"
	"io"
	"net"
	"os"
	"time"
)

func newMessage() ioadapter.IMessage {
	msg := message.NewUploadSubscriberPacketMessage()
	nfInstID := &element.NFInstanceIDElement{}
	nfInstID.Value = []byte("f2b43e11-762a-49b9-b5d4-4c3452ecb8f9")
	msg.NFInstanceIDElement = nfInstID
	logger.DebugFmt("New Message, nfInstID value: %v, len: %d", string(nfInstID.Value), len(nfInstID.Value))

	subIdentities := &element.SubscriberIdentitiesElement{}
	subIdentities.IMSIElement = &element.IMSIElement{}
	subIdentities.IMSIElement.Value = []byte("1111111111")
	//subIdentities.Flag = subIdentities.Flag | 1

	subIdentities.IMEIElement = &element.IMEIElement{}
	subIdentities.IMEIElement.Value = []byte("2222222222")
	//subIdentities.Flag = subIdentities.Flag | (1 << 2)

	subIdentities.UEIPV4Element = &element.UEIPV4Element{}
	subIdentities.UEIPV4Element.IP = net.ParseIP("172.10.0.15")
	//subIdentities.Flag = subIdentities.Flag | (1 << 4)
	msg.SubscriberIdentitiesElement = subIdentities

	extFieldElement := &element.ExtensionFieldElement{
		ProtocolType: 2, MicroserviceType: 3,
	}
	extFieldElement.Flag = 0x2
	msg.ExtensionFieldElement = extFieldElement

	subPacket := &element.PayloadPacketElement{}
	subPacket.Value = []byte("ccc-vvv-bbb-nnn-mmm")
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

	layerPkt := ioadapter.NewLayerPacket()
	layerPkt.MsgType = 1
	layerPkt.Msg = msg

	writer, err := newWriter(*addr)
	if err != nil {
		logger.ErrorFmt("NewWriter Failed, %v", err)
		os.Exit(1)
	}

	buf := bufferpool.AllocateBufferObject()
	if err := layerPkt.MarshalBinary(buf); err != nil {
		logger.ErrorFmt("MarshalBinary, %v", err)
	}

	writer.Write(buf.Bytes())
	time.Sleep(time.Second * 2)
	os.Exit(0)
}
