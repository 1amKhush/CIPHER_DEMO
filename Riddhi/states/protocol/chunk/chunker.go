package chunk

import (
 "io"
 "os"
 "sync"
 "sort"
 "runtime"
 
 "riddhi/states/protocol/crypto"

)

// ChunklargeFile splits a large file into smaller chunks in parallel and returns metadata for each chunk.
// It divides the file into chunks and processes them concurrently using multiple goroutines.
func (c *DefaultFileChunker) ChunklargeFile(filePath string) (*ChunkedFile, error) {

 var wg sync.WaitGroup

 //used to protect shared access to the chunks slice-but we will use channel to avoid mutex
 //var mu sync.Mutex 

 // Store metadata for each chunk
 var chunks []ChunkMeta

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

 // Calculate the number of chunks based on the file size and chunk size
 numChunks := int(fileInfo.Size() / int64(c.ChunkSize))
 if fileInfo.Size()%int64(c.ChunkSize) != 0 {
  numChunks++
 }

 // Create channels to communicate between goroutines
 errChan := make(chan error, numChunks)
 resultChan := make(chan ChunkMeta, numChunks)
 jobChan := make(chan ChunkJob, numChunks)

 // Populate the job channel with chunk indices
for i := 0; i < numChunks; i++ {

    offset := int64(i) * int64(c.ChunkSize)

    jobChan <- ChunkJob{
        Index:  i,
        Offset: offset,
    }
}

close(jobChan)

//create fileId(raw binary bytes)
fullData, err := os.ReadFile(filePath)
if err != nil {
	return nil, err
}
fileID := crypto.GenerateHash(fullData)


 // Start multiple goroutines to process chunks in parallel
 //instead of hardcoded worker,choose worker count depending on machine
 workers := runtime.NumCPU()

for i := 0; i < workers; i++ { // Number of parallel workers
  wg.Add(1)
  go func() {
   defer wg.Done()
   for job := range jobChan {
   
    // Create a buffer for chunk data
    buffer := make([]byte, c.ChunkSize) 

    // Read chunkSize bytes from the file into the buffer
    bytesRead, err := file.ReadAt(buffer, job.Offset)
    if err != nil && err != io.EOF {
     errChan <- err
     return
    }

    // If bytesRead is 0, it means EOF is reached
    if bytesRead > 0 {

    chunkCopy := make([]byte, bytesRead)
    copy(chunkCopy, buffer[:bytesRead])

 //generate hash
 chunkHash :=
	crypto.GenerateHash(chunkCopy)

  //genrate leaf hash
  leafHash := crypto.GenerateLeafHash(
    fileID,
    job.Index,
    bytesRead,
    chunkCopy,
  )

 //MetaData creation
 meta := ChunkMeta{
    ChunkIndex: job.Index,
    ChunkSize: bytesRead,
    ChunkData:  chunkCopy,
    ChunkHash:  chunkHash,
    LeafHash:   leafHash,
 }


    // Send the metadata to the result channel
    resultChan <- meta 

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

//to maintain order sort the chunks based on ChunkIndex because goroutines may finish in different order
sort.Slice(chunks, func(i, j int) bool {
    return chunks[i].ChunkIndex <
           chunks[j].ChunkIndex
}) 

 // store the leaf hash of each chunk
leafHashes := make([][32]byte, 0, len(chunks))
for _, chunk := range chunks {
    leafHashes = append(leafHashes, chunk.LeafHash)
}

root := crypto.BuildMerkleRoot(leafHashes)


 return &ChunkedFile{
     FileID: fileID,
     FileSize:fileInfo.Size(),
     ChunkSize:c.ChunkSize,
     Chunks: chunks,
     LeafHashes: leafHashes, 
     MerkleRoot: root, 
       }, nil
}