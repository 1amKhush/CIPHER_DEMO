package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// acts as seeder to cut file into chunks, encrypt them, save them, and serve to leechers

// take the file pointer determine the number of chunks and max of 100 chunks ans min size of 38 byte  

// take chunk data genrate the 32 byte key and 16 byte nonce for each chunk using rand and encrypt the data usign aes gcm 

func storeChunkData(chunkIndex int, data []byte, keys KeyNonceMap) []byte {
	// Generate a secure random 32-byte key for each chunk
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	block2, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, block2.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		panic(err)
	}

	ciphertext := block2.Seal(nil, nonce, data, nil)

	keys[chunkIndex] = KeyNoncePair{Key: key, Nonce: nonce}

	return ciphertext
}

func splitAndEncryptChunks(data []byte, chunkSize int, keys KeyNonceMap) [][]byte {
	var chunks [][]byte
	chunkIndex := 0
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		
		ciphertext := storeChunkData(chunkIndex, data[i:end], keys)
		chunks = append(chunks, ciphertext)
		chunkIndex++
	}
	return chunks
}

func PreprocessFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	stat, _ := file.Stat()
	fmt.Println("File Size:", stat.Size())
	buffer := make([]byte, stat.Size())
	file.Read(buffer)

	chunkSize := int(stat.Size() / 100)
	if chunkSize < 38 {
		chunkSize = 38
	}
	fmt.Println("Dynamic Chunk Size:", chunkSize)

	keys := make(KeyNonceMap)
	fmt.Println("Chunking and encrypting...")
	chunks := splitAndEncryptChunks(buffer, chunkSize, keys)

	fmt.Println("Generating Merkle tree...")
	merkleTree := generateMerkleTree(chunks)

	metadata := Metadata{
		FileName:       filepath.Base(filename),
		Extension:      filepath.Ext(filename),
		ChunkCount:     len(chunks),
		MerkleRootHash: fmt.Sprintf("%x", merkleTree.RootNode.Hash),
	}

	// Save to disk
	baseName := filepath.Base(filename)
	outDir := filepath.Join("chunks", baseName)
	os.MkdirAll(outDir, 0755)

	// Save chunks
	for i, chunk := range chunks {
		chunkFile := filepath.Join(outDir, fmt.Sprintf("%s_c%d.bin", baseName, i))
		os.WriteFile(chunkFile, chunk, 0644)
	}

	// Save metadata
	metaBytes, _ := json.MarshalIndent(metadata, "", "  ")
	os.WriteFile(filepath.Join(outDir, "metadata.json"), metaBytes, 0644)

	// Save keys
	keysBytes, _ := json.MarshalIndent(keys, "", "  ")
	os.WriteFile(filepath.Join(outDir, "keys.json"), keysBytes, 0644)

	fmt.Println("Preprocessed and stored to disk in", outDir)
}

func RunSeeder(trackerAddr, filename, myIP, listenPort string) {
	// First ensure file is preprocessed
	PreprocessFile(filename)
	
	listenAddr := myIP + ":" + listenPort

	// Announce to Tracker
	fmt.Println("Announcing to tracker:", trackerAddr)
	if err := AnnounceToTracker(trackerAddr, filepath.Base(filename), listenAddr); err != nil {
		fmt.Println("Failed to announce to tracker:", err)
		return
	}
	fmt.Println("Successfully announced to tracker.")

	var t Transport = &TCPTransport{}
	
	handler := func(reqFilename string) (Metadata, KeyNonceMap, [][]byte, error) {
		baseName := filepath.Base(reqFilename)
		outDir := filepath.Join("chunks", baseName)

		metaBytes, err := os.ReadFile(filepath.Join(outDir, "metadata.json"))
		if err != nil { return Metadata{}, nil, nil, err }
		var metadata Metadata
		json.Unmarshal(metaBytes, &metadata)

		keysBytes, err := os.ReadFile(filepath.Join(outDir, "keys.json"))
		if err != nil { return Metadata{}, nil, nil, err }
		var keys KeyNonceMap
		json.Unmarshal(keysBytes, &keys)

		var chunks [][]byte
		for i := 0; i < metadata.ChunkCount; i++ {
			chunkFile := filepath.Join(outDir, fmt.Sprintf("%s_c%d.bin", baseName, i))
			chunkBytes, err := os.ReadFile(chunkFile)
			if err != nil { return Metadata{}, nil, nil, err }
			chunks = append(chunks, chunkBytes)
		}

		return metadata, keys, chunks, nil
	}

	if err := t.ListenForRequests(listenPort, handler); err != nil {
		fmt.Println("Seeder server error:", err)
	}
}
