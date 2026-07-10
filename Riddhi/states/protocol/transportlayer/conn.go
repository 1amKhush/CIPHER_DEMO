package transportlayer

import (
	//"net"

	"github.com/libp2p/go-libp2p/core/network"
	"riddhi/states/protocol/packetlayer"
)


//for p2p
type Transport struct {
    Stream network.Stream  //instead of conn net.Conn
}

//SendPacket
func (t *Transport) SendPacket(
	packet *packetlayer.Packet,
) error {

	encoded := packetlayer.EncodePacket(
		byte(packet.Type),
		packet.Payload,
	)

	_, err := t.Stream.Write(encoded)

	return err
}

//ReceivePacket
func (t *Transport) ReceivePacket() (
	*packetlayer.Packet,
	error,
) {
	return packetlayer.DecodePacket(
		t.Stream,
	)
}



