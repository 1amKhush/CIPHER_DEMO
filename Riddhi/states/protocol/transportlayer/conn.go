package transportlayer

import (
	"net"
	"riddhi/states/protocol/packetlayer"
)

type Transport struct {
	Conn net.Conn
}

//SendPacket
func (t *Transport) SendPacket(
	packet *packetlayer.Packet,
) error {

	encoded := packetlayer.EncodePacket(
		byte(packet.Type),
		packet.Payload,
	)

	_, err := t.Conn.Write(encoded)

	return err
}

//ReceivePacket
func (t *Transport) ReceivePacket() (
	*packetlayer.Packet,
	error,
) {
	return packetlayer.DecodePacket(
		t.Conn,
	)
}



