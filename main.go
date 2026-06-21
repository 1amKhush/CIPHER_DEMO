package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)
//creating job queue
type Job struct {
	Index int
	Chunk []byte
}

type Result struct {
	Info ChunkInfo
}
//assigning worker for each job queue
func worker(
	id int,
	jobs <-chan Job,
	results chan<- Result,
	fileID []byte,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for job := range jobs {

		key, commitment, err := GenerateCommitment(job.Chunk)
		if err != nil {
			panic(err)
		}

		leaf := GenerateLeaf(
			fileID,
			job.Index,
			job.Chunk,
		)

		info := ChunkInfo{
			Index:      job.Index,
			Size:       len(job.Chunk),
			Key:        key,
			Data:       job.Chunk,
			Commitment: commitment,
			Leaf:       leaf,
		}

		results <- Result{Info: info}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("go run . <file>")
		return
	}

	filename := os.Args[1]

	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	fileHash := sha256.Sum256(data)
	fileID := fileHash[:]

	chunks := SplitIntoChunks(data)

	numWorkers := runtime.NumCPU()
	//time measurement
	processStart := time.Now()
    //making channels for goroutines
	jobs := make(chan Job, len(chunks))
	results := make(chan Result, len(chunks))

	var wg sync.WaitGroup

	// Start workers
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, fileID, &wg)
	}

	// Send jobs
	for i, chunk := range chunks {
		jobs <- Job{
			Index: i,
			Chunk: chunk,
		}
	}
	close(jobs)

	// Wait for workers
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	resultMap := make(map[int]ChunkInfo)

	for r := range results {
		resultMap[r.Info.Index] = r.Info
	}

	fmt.Printf("Generated all hashes in: %s\n", time.Since(processStart))
}
