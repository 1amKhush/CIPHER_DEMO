package packetlayer

type KeyRevealPayload struct {
	ChunkIndex uint32
	Key         [32]byte
}

type TicketRejectPayload struct {
	Reason uint8
}