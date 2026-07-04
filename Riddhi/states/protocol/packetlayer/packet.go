package packetlayer

type Packet struct {
	Type    PacketType
	Length  uint32
	Payload []byte
}