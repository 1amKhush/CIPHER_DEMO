package session

import(
    "crypto/ed25519"
    "riddhi/states/protocol/crypto"
)
// session holds temporary runtime state
// for an active peer-to-peer transfer

type Session struct {
    // Requester identity
    RequesterPrivateKey ed25519.PrivateKey 
    RequesterPublicKey ed25519.PublicKey 

    // provider identity
    ProviderPrivateKey ed25519.PrivateKey 
    ProviderPublicKey ed25519.PublicKey

    CurrentChunkIndex uint32
    StoredCiphertext  []byte
    StoredCommitment  [32]byte
    StoredNonce       [24]byte
    StoredKey         [32]byte
    OriginalChunk     []byte
    
    FileID [32]byte
    
    Proof      crypto.MerkleProof
    MerkleRoot        [32]byte
}

