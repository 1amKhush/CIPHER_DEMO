package main

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"golang.org/x/crypto/sha3"
)

// ChunkInfo stores all information related to a chunk.
type ChunkInfo struct {
	Index      int
	Data       []byte
	Key        []byte
	Commitment []byte
	Leaf       []byte
}

// Metadata stored inside metadata.json for each chunk.
type ChunkMeta struct {
	Index      int    `json:"index"`
	Commitment string `json:"commitment"`
	Leaf       string `json:"leaf"`
}

// File metadata stored alongside chunk files.
type FileMeta struct {
	FileName    string      `json:"file_name"`
	FileID      string      `json:"file_id"`
	ChunkSize   int         `json:"chunk_size"`
	TotalChunks int         `json:"total_chunks"`
	Chunks      []ChunkMeta `json:"chunks"`
}

// For best latency it should be ~1200 Bytes for QUIC protocol for live streaming
// Fixed chunk size specified by the protocol (32 KB).
const ChunkSize = 32768 // 32 KB

// Converts a uint64 value into an 8-byte Big Endian representation.
// Used when serializing ChunkIndex and Length for leaf generation.
func Uint64ToBytes(v uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, v)

	return buf
}

// Computes the Keccak256 hash of arbitrary data.
// Used for FileID, Commitments and Merkle Leaves.
func Keccak256(data []byte) []byte {
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(data)

	return hasher.Sum(nil)
}

// Generates a cryptographically secure 32-byte random key.
// Each chunk receives its own random key.

func GenerateRandomKey() ([]byte, error) {
	key := make([]byte, 32)

	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// Commitment = Keccak256(Key || ChunkData)
//
// This acts as a cryptographic commitment to the chunk.
func GenerateCommitment(key []byte, chunk []byte) []byte {
	input := make([]byte, 0, len(key)+len(chunk))

	input = append(input, key...)
	input = append(input, chunk...)

	return Keccak256(input)
}

// FileID = Keccak256(EntireFile)
//
// Uniquely identifies a file regardless of filename.
func GenerateFileID(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return Keccak256(data), nil
}

// Leaf = Keccak256(FileID || ChunkIndex || Length || ChunkData)
// Used as the leaf node when constructing a Merkle Tree.

func GenerateLeaf(fileID []byte, index int, chunk []byte) []byte {

	indexBytes := Uint64ToBytes(uint64(index))
	lengthBytes := Uint64ToBytes(uint64(len(chunk)))

	input := make([]byte, 0, len(fileID)+8+8+len(chunk))

	input = append(input, fileID...)
	input = append(input, indexBytes...)
	input = append(input, lengthBytes...)
	input = append(input, chunk...)

	return Keccak256(input)
}

// Sequential implementation.
//
// Reads the file chunk-by-chunk and immediately generates:
//   - Random Key
//   - Commitment
//   - Merkle Leaf
//
// Returns a slice containing all processed chunks.
// Process time is ~16ms
func ProcessFile(path string, fileID []byte) ([]ChunkInfo, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var result []ChunkInfo
	index := 0

	for {
		buffer := make([]byte, ChunkSize)

		n, err := file.Read(buffer)

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		buffer = buffer[:n]

		key, err := GenerateRandomKey()
		if err != nil {
			return nil, err
		}

		result = append(result, ChunkInfo{
			Index:      index,
			Data:       buffer,
			Key:        key,
			Commitment: GenerateCommitment(key, buffer),
			Leaf:       GenerateLeaf(fileID, index, buffer),
		})

		index++
	}

	return result, nil
}

// Parallel implementation.

// Reads the file and returns raw chunks.
//
// Intended to be used with ProcessChunksParallel().

func ReadChunks(path string) ([][]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var chunks [][]byte

	for {
		buffer := make([]byte, ChunkSize)
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		chunks = append(chunks, buffer[:n])
	}

	return chunks, nil
}

// Each chunk is processed independently inside its own goroutine.
// Chunk ordering is preserved using the chunk index.

// Process time is ~13ms
func ProcessChunksParallel(chunks [][]byte, fileID []byte) ([]ChunkInfo, error) {
	result := make([]ChunkInfo, len(chunks))
	var wg sync.WaitGroup

	errChan := make(
		chan error,
		len(chunks),
	)

	for i, chunk := range chunks {
		wg.Add(1)

		go func(index int, data []byte) {
			defer wg.Done()

			key, err := GenerateRandomKey()
			if err != nil {
				errChan <- err
				return
			}

			// Store result at the correct index.
			// Safe because each goroutine writes to a unique position.
			result[index] = ChunkInfo{
				Index:      index,
				Data:       data,
				Key:        key,
				Commitment: GenerateCommitment(key, data),
				Leaf:       GenerateLeaf(fileID, index, data),
			}

		}(i, chunk)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// Persists all chunk files and metadata to disk.
//
// Directory Structure:
//
// storage/<FileID>/
// ├── 0.chunk
// ├── 1.chunk
// ├── ...
// └── metadata.json

func SaveFile(fileID []byte, path string, chunks []ChunkInfo) error {

	dir := fmt.Sprintf("storage/%x", fileID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	meta := FileMeta{
		FileName:    path,
		FileID:      fmt.Sprintf("%x", fileID),
		ChunkSize:   ChunkSize,
		TotalChunks: len(chunks),
	}

	for _, chunk := range chunks {

		filename := fmt.Sprintf("%s/%d.chunk", dir, chunk.Index)

		err := os.WriteFile(filename, chunk.Data, 0644)
		if err != nil {
			return err
		}

		meta.Chunks = append(meta.Chunks, ChunkMeta{
			Index:      chunk.Index,
			Commitment: fmt.Sprintf("%x", chunk.Commitment),
			Leaf:       fmt.Sprintf("%x", chunk.Leaf),
		},
		)
	}

	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(dir+"/metadata.json", data, 0644)
}

func main() {
	// Measure total execution time.
	start := time.Now()
	// Input file to process.
	path := "CSE.pdf"

	// Generate a unique identifier for the file.
	fileID, err := GenerateFileID(path)
	if err != nil {
		panic(err)
	}

	// Alternate Sequential processing pipeline.
	// infos, err := ProcessFile(path, fileID)
	// if err != nil {
	// 	panic(err)
	// }

	chunks, err := ReadChunks(path)
	if err != nil {
		panic(err)
	}

	// Parallel pipeline:
	//
	// chunks, err := ReadChunks(path)
	// infos, err := ProcessChunksParallel(chunks, fileID)
	infos, err := ProcessChunksParallel(chunks, fileID)
	if err != nil {
		panic(err)
	}

	// Persist chunks and metadata.
	err = SaveFile(fileID, path, infos)
	if err != nil {
		panic(err)
	}

	// Display execution statistics.
	fmt.Println("Done")
	fmt.Printf("Processed %d chunks in %v\n", len(infos), time.Since(start))
}
