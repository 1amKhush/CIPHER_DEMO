package packetlayer

import (
	"crypto/ed25519"
)
type LotteryTicketPayload struct {
	ChunkIndex uint32
	TicketID   uint32
	Signature  [64]byte
	ClientPublicKey ed25519.PublicKey
}

