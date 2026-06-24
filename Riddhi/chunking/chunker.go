package main

import (
 //"encoding/hex"
 //"fmt"
 "io"
 "os"
 "sync"
 "sort"
 "runtime"

// "golang.org/x/crypto/sha3"
)

// ChunkFile splits a file into smaller chunks and returns metadata for each chunk.
// It reads the file sequentially and chunks it based on the specified chunk size.
func (c *DefaultFileChunker) ChunkFile(filePath string) ([]ChunkMeta, error) {
 var chunks []ChunkMeta // Store metadata for each chunk

 // Open the file for reading
 file, err := os.Open(filePath)
 if err != nil {
  return nil, err
 }
 defer file.Close()

 //Creation of fileid
 fileID := GenerateHash([]byte(filePath))

 // Create a buffer to hold the chunk data
 buffer := make([]byte, c.chunkSize)
 index := 0 // Initialize chunk index

 // Loop until EOF is reached
 for {
  // Read chunkSize bytes from the file into the buffer
  bytesRead, err := file.Read(buffer)
  if err != nil && err != io.EOF {
   return nil, err
  }
  if bytesRead == 0 {
   break // If bytesRead is 0, it means EOF is reached
  }

//   // Construct the chunk file name
//   chunkFileName := fmt.Sprintf("%s.chunk.%d", filePath, index)

//   // Create a new chunk file and write the buffer data to it
//   chunkFile, err := os.Create(chunkFileName)
//   if err != nil {
//    return nil, err
//   }

  chunkCopy := make([]byte, bytesRead)
 copy(chunkCopy, buffer[:bytesRead])

//   _, err = chunkFile.Write(chunkCopy)
//   if err != nil {
//    return nil, err
//   }

//   // Close the chunk file
//   chunkFile.Close()

  //commitment engine
  key, commitment, err := GenerateCommitment(chunkCopy)
  if err != nil {
    return nil, err
   }

  //generate Leaf

  leaf := GenerateLeafHash(
    fileID,
    index,
    bytesRead,
    chunkCopy,
)

//Append metadata
  chunks = append(chunks, ChunkMeta{
   // FileName:       chunkFileName,
    ChunkIndex:          index,
    ChunkSize:      bytesRead,
    RandomKey:      key,
    CommitmentHash: commitment,
    MerkleLeafHash: leaf,
})

  // Move to the next chunk
  index++
 }

 return chunks, nil
}


//for future
// ChunklargeFile splits a large file into smaller chunks in parallel and returns metadata for each chunk.
// It divides the file into chunks and processes them concurrently using multiple goroutines.
func (c *DefaultFileChunker) ChunklargeFile(filePath string) ([]ChunkMeta, error) {
 var wg sync.WaitGroup
 //var mu sync.Mutex //use later
 var chunks []ChunkMeta // Store metadata for each chunk

 // Open the file for reading
 file, err := os.Open(filePath)
 if err != nil {
  return nil, err
 }
 defer file.Close()

 // Get file information to determine the number of chunks
 fileInfo, err := file.Stat()
 if err != nil {
  return nil, err
 }

 numChunks := int(fileInfo.Size() / int64(c.chunkSize))
 if fileInfo.Size()%int64(c.chunkSize) != 0 {
  numChunks++
 }

 // Create channels to communicate between goroutines
 errChan := make(chan error, numChunks)
 resultChan := make(chan ChunkMeta, numChunks)
 jobChan := make(chan ChunkJob, numChunks)

 // Populate the job channel with chunk indices
for i := 0; i < numChunks; i++ {

    offset := int64(i) * int64(c.chunkSize)

    jobChan <- ChunkJob{
        Index:  i,
        Offset: offset,
    }
}

close(jobChan)

//create fileId
fileID := GenerateHash([]byte(filePath))

 // Start multiple goroutines to process chunks in parallel
 //instead of hardcoded worker,choose worker count depending on machine
 workers := runtime.NumCPU()

 for i := 0; i < workers; i++ { // Number of parallel workers
  wg.Add(1)
  go func() {
   defer wg.Done()
   for job := range jobChan {
   
    // Create a buffer for chunk data
    buffer := make([]byte, c.chunkSize) 

    // Read chunkSize bytes from the file into the buffer
    bytesRead, err := file.ReadAt(buffer, job.Offset)
    if err != nil && err != io.EOF {
     errChan <- err
     return
    }

    // If bytesRead is 0, it means EOF is reached
    if bytesRead > 0 {

   
    //  // Construct the chunk file name
    //  chunkFileName := fmt.Sprintf("%s.chunk.%d", filePath, job.Index)

    //  // Create a new chunk file and write the buffer data to it
    //  chunkFile, err := os.Create(chunkFileName)
    //  if err != nil {
    //   errChan <- err
    //   return
    //  }

    chunkCopy := make([]byte, bytesRead)
    copy(chunkCopy, buffer[:bytesRead])

    //  _, err = chunkFile.Write(chunkCopy)
    //  if err != nil {
    //   errChan <- err
    //   return
    //  }

    //Generate Hash
     key, commitment, err := GenerateCommitment(chunkCopy)
     if err != nil {
     errChan <- err
     return
     }

    //Generate Leaf
    leaf := GenerateLeafHash(
    fileID,
    job.Index,
    bytesRead,
    chunkCopy,
    )


  //MetaData creation
    meta := ChunkMeta{
   // FileName: chunkFileName,
    ChunkIndex: job.Index,
    ChunkSize: bytesRead,

    RandomKey: key,

    CommitmentHash: commitment,

    MerkleLeafHash: leaf,
}

resultChan <- meta 

   //not needed
    //  // Append metadata for the chunk to the chunks slice
    //  chunk := ChunkMeta{
    //   FileName: chunkFileName,ChunkIndex: index,ChunkSize: bytesRead,
    //  }
    //  mu.Lock()
    //  chunks = append(chunks, chunk)
    //  mu.Unlock()

     // Close the chunk file
    // chunkFile.Close()

    }
   }
  }()
 }

 // Wait for all goroutines to finish
 go func() {
  wg.Wait()
  
  close(errChan)
  close(resultChan)
 }()

for meta := range resultChan {
    chunks = append(chunks, meta)
}
select {
case err := <-errChan:
    if err != nil {
        return nil, err
    }
default:
}
//to maintain order
sort.Slice(chunks, func(i, j int) bool {
    return chunks[i].ChunkIndex <
           chunks[j].ChunkIndex
}) 

 return chunks, nil
}