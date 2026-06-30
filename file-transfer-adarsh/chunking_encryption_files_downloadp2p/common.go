package main

import (
	"golang.org/x/crypto/sha3"
)

// holds shared stuff like tree structures and hash maker used by everyone

type MerkleTree struct {
	RootNode *Node
}

type Metadata struct {
	FileName       string `json:"file_name"`
	Extension      string `json:"extension"`
	ChunkCount     int    `json:"chunk_count"`
	MerkleRootHash string `json:"merkle_root_hash"`
}

type Node struct {
	Hash  []byte
	Left  *Node
	Right *Node
}

func calculateHash(data []byte) []byte {
	h := sha3.NewLegacyKeccak256()
	h.Write(data)
	return h.Sum(nil)
}

func buildTree(leaves []*Node) *Node {
	if len(leaves) == 1 {
		return leaves[0]
	}
	var newLevel []*Node
	for i := 0; i < len(leaves)-1; i += 2 {
		left := leaves[i]
		right := leaves[i+1]
		newLevel = append(newLevel, &Node{Hash: calculateHash(append(left.Hash, right.Hash...)), Left: left, Right: right})
	}

	if len(leaves)%2 != 0 {
		newLevel = append(newLevel, leaves[len(leaves)-1])
	}
	return buildTree(newLevel)
}

func generateMerkleTree(chunks [][]byte) *MerkleTree {
	var leaves []*Node
	for _, chunk := range chunks {
		leaves = append(leaves, &Node{Hash: calculateHash(chunk)})
	}
	return &MerkleTree{RootNode: buildTree(leaves)}
}

type KeyNoncePair struct {
	Key   []byte `json:"key"`
	Nonce []byte `json:"nonce"`
}

// Map from chunk index to KeyNoncePair
type KeyNonceMap map[int]KeyNoncePair
