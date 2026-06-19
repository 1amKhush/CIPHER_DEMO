package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"

	"golang.org/x/crypto/sha3"
)

const chunkSize = 32768
const threshold int64 = 1 * 1024 * 1024 * 1024

func ChunkFileS(filePath string, chunkSize int64) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()
	buffer := make([]byte, int(chunkSize))
	chunkIndex := 0
	fileId := make([]byte, 32)
	_, newererr := rand.Read(fileId)
	if newererr != nil {
		return
	}
	for {
		bytes, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Println(err)
			return
		}
		if bytes == 0 {
			break
		}
		chunkCont := buffer[:bytes]
		key := make([]byte, 32)
		_, newerr := rand.Read(key)
		if newerr != nil {
			fmt.Println(newerr)
			return
		}
		fmt.Printf("Chunk %d\n", chunkIndex)
		byteIndex := make([]byte, 8)
		binary.BigEndian.PutUint64(byteIndex, uint64(chunkIndex))
		byteLength := make([]byte, 8)
		binary.BigEndian.PutUint64(byteLength, uint64(bytes))
		hash := sha3.NewLegacyKeccak256()
		hash.Write(key)
		hash.Write(chunkCont)
		H_resp := hash.Sum(nil)
		fmt.Printf("Commitment hash: %x\n", H_resp)

		newHash := sha3.NewLegacyKeccak256()
		newHash.Write(fileId)
		newHash.Write(byteIndex)
		newHash.Write(byteLength)
		newHash.Write(chunkCont)
		Leaf := newHash.Sum(nil)
		fmt.Printf("Merkle Leaf: %x\n", Leaf)

		fmt.Printf("Chunk Data: %x\n", chunkCont)

		chunkIndex++
		if err == io.EOF {
			break
		}
	}
}

func ChunkFileP(filePath string, chunkSize int64, workersNum int) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return
	}
	fileSize := fileInfo.Size()
	totalChunks := (fileSize + chunkSize - 1) / chunkSize
	chunkIndex := make(chan int64, totalChunks)
	for i := int64(0); i < totalChunks; i++ {
		chunkIndex <- i
	}
	close(chunkIndex)
	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	for w := 0; w < workersNum; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range chunkIndex {
				select {
				case err := <-errChan:
					if err != nil {
						fmt.Println(err)
					}
					return
				default:
				}
				offset := idx * chunkSize
				length := int64(chunkSize)
				if offset+length > fileSize {
					length = fileSize - offset
				}
				chunkData := make([]byte, int(length))
				_, err := file.ReadAt(chunkData, offset)
				if err != nil && err != io.EOF {
					select {
					case errChan <- fmt.Errorf("error at %d", offset):
					default:
					}
					return
				}

			}
		}()

	}

	wg.Wait()
	close(errChan)
	if err := <-errChan; err != nil {
		return
	}

}
func main() {
	fmt.Println("Enter file's path:")
	var path string
	_, err := fmt.Scanln(&path)
	if err != nil {
		fmt.Println(err)
		return
	}
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("There was an error opening your file as:", err)
		return
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	fileSize := fileInfo.Size()
	if fileSize <= threshold {
		ChunkFileS(path, chunkSize)
	} else {
		numWorkers := runtime.NumCPU()
		ChunkFileP(path, chunkSize, numWorkers)
	}
}
