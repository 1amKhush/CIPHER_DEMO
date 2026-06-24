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

//threshold is the file size (1 GB) above which we switch to parallel processing.
const threshold int64 = 1 * 1024 * 1024 * 1024  

/* ChunkFileS reads a file sequentially and for each chunk creates:
 a random key
 a commitment hash 
 a Merkle leaf hash */

func ChunkFileS(filePath string, chunkSize int64) {
	file, err := os.Open(filePath)     
	if err != nil {
		return
	}
	defer file.Close()
  
// Creating a buffer to hold chunk

	buffer := make([]byte, int(chunkSize))
	chunkIndex := 0
	fileId := make([]byte, 32)
	_, newererr := rand.Read(fileId)
	if newererr != nil {
		return
	}
  
// Read up to designated chunkSize bytes into our buffer

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

    /* ChunkIndex and bytes are stored as integers but cryptographic hash functions only accept a raw stream of bytes 
    to store them allocating two 8 bytes slices for Chunk Index and bytes using Big-Endian byte ordering (means the most significant bytes come first) */
    
    //Don't convert into strings as it might result into hash collisions
    
		byteIndex := make([]byte, 8)
		binary.BigEndian.PutUint64(byteIndex, uint64(chunkIndex))
		byteLength := make([]byte, 8)
		binary.BigEndian.PutUint64(byteLength, uint64(bytes))

    // Creation of Commitment hash 
		// Keccak256( key || chunk data )

		hash := sha3.NewLegacyKeccak256()
		hash.Write(key)
		hash.Write(chunkCont)
		H_resp := hash.Sum(nil)
		fmt.Printf("Commitment hash: %x\n", H_resp)

    // Merkle leaf hash 
		// Keccak256( fileId || index || length || chunk data )


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

// ChunkFileP does file chunking using multiple goroutines 

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
  
  // Used (chunkSize - 1) to round up any partial chunk

  totalChunks := (fileSize + chunkSize - 1) / chunkSize

  // Creating a buffered channel that will hold all the chunk indices

	chunkIndex := make(chan int64, totalChunks)
	for i := int64(0); i < totalChunks; i++ {     
		chunkIndex <- i
	}
	close(chunkIndex)
	var wg sync.WaitGroup   // To wait for all workers to finish

	errChan := make(chan error, 1)

  // Starting 'workersNum' goroutines
  for w := 0; w < workersNum; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
      // Each worker will take one chunk's index out of the channel and process that chunk
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
          //To make sure that when length is not 32kbs it still gets chunked in its own size only
					length = fileSize - offset
				}
				chunkData := make([]byte, int(length))
        // used file.ReadAt as it is preferrable for concurrent use 
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
  //Ask user for file's path
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
		numWorkers := runtime.NumCPU()   // important as it tells to use as many workers as CPU cores for maximised output

		ChunkFileP(path, chunkSize, numWorkers)
	}
}
