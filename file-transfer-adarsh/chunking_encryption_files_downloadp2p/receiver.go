package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"
	"time"
)

// acts as leecher to ask tracker, download chunks, check them, and join them
func combineChunks(metadata Metadata, keys KeyNonceMap, chunks [][]byte) {
	fmt.Println("Combining chunks for file:", metadata.FileName)
	
	// Ensure an output directory exists
	os.MkdirAll("received_files", 0755)
	
	outputFile := "received_files/output_" + metadata.FileName

	f, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for i := 0; i < metadata.ChunkCount; i++ {
		keyNonce := keys[i]

		block, err := aes.NewCipher(keyNonce.Key)
		if err != nil {
			panic(err)
		}

		block2, err := cipher.NewGCM(block)
		if err != nil {
			panic(err)
		}

		ciphertext := chunks[i]

		decryptedGCM, err := block2.Open(nil, keyNonce.Nonce, ciphertext, nil)
		if err != nil {
			panic(err)
		}

		if _, err := f.Write(decryptedGCM); err != nil {
			panic(err)
		}
	}
	
	fmt.Println("File reconstructed successfully at", outputFile)
}

func RunLeecher(trackerAddr string, filename string) {
	fmt.Println("Querying tracker", trackerAddr, "for file:", filename)
	peerAddr, err := QueryTracker(trackerAddr, filename)
	if err != nil {
		fmt.Println("Tracker query failed:", err)
		return
	}
	fmt.Println("Tracker responded with peer address:", peerAddr)

	var t Transport = &TCPTransport{}
	fmt.Println("Connecting to peer", peerAddr, "to download...")
	
	metadata, keys, chunks, err := t.RequestFile(peerAddr, filename)
	if err != nil {
		fmt.Println("Error receiving:", err)
		return
	}

	fmt.Println("Received Metadata:")
	fmt.Println(" - FileName:", metadata.FileName)
	fmt.Println(" - Chunks:", metadata.ChunkCount)
	fmt.Println(" - MerkleRoot:", metadata.MerkleRootHash)

	start := time.Now()
	// Re-verify Merkle Tree
	fmt.Println("Verifying Merkle Tree...")
	merkleTree := generateMerkleTree(chunks)
	calculatedRoot := fmt.Sprintf("%x", merkleTree.RootNode.Hash)
	if calculatedRoot != metadata.MerkleRootHash {
		fmt.Println("CRITICAL ERROR: Merkle tree root hash mismatch!")
		fmt.Println("Expected:", metadata.MerkleRootHash)
		fmt.Println("Calculated:", calculatedRoot)
		return
	}
	fmt.Println("Merkle tree verification passed.")

	combineChunks(metadata, keys, chunks)
	fmt.Printf("Reconstruction took: %v\n", time.Since(start))
}
