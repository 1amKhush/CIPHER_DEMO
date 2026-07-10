package packetlayer

import (
	"riddhi/states/protocol/crypto"
	"crypto/ed25519"
)

type FileMetadataResponsePayload struct {
	FileID    [32]byte
	FileSize  uint64
	ChunkSize uint32
}

type ChunkResponsePayload struct {
	FileID [32]byte
	ChunkIndex uint32
	Nonce      [24]byte
	Ciphertext []byte
	Commitment [32]byte
	Proof      crypto.MerkleProof
	MerkleRoot        [32]byte
	ProviderSignature [64]byte
	ProviderPublicKey ed25519.PublicKey
}