package main

import (
	"fmt"
	"log"
	//"os"
)

func main() {
    fmt.Println("--Sequential Chunking--")
    chunker := DefaultFileChunker{
        chunkSize: 32768,
    }

    chunks, err := chunker.ChunkFile("test.pdf")
    if err != nil {
        panic(err)
    }

   for _, chunk := range chunks {
    fmt.Printf(
        "Chunk %d | Size: %d\n",
        chunk.ChunkIndex,
        chunk.ChunkSize,
       
    )
}
//use if chunks are stored
// for _, chunk := range chunks {

//     data, err := os.ReadFile(chunk.FileName)
//     if err != nil {
//         panic(err)
//     }

//     hash := GenerateHash(data)

//     fmt.Printf(
//         "Chunk %d | Hash: %s\n",
//         chunk.ChunkIndex,
//         hash,
//     )
// }
//Commitment test
data := []byte("hello world")
//data1:=[]byte("hello worlds")

key, commitment, err := GenerateCommitment(data)
if err!=nil{
    log.Fatal(err)
}

valid := VerifyCommitment(
    key,
    data,
    commitment,
)

fmt.Println(valid)

//leaf test
leaf := GenerateLeafHash(
    "file123",
    0,
    len(data),
    data,
)

fmt.Println(leaf)


//Sequential metadata

for _, chunk := range chunks {

    fmt.Println("Chunk:", chunk.ChunkIndex)

    fmt.Println("Commitment:", chunk.CommitmentHash)

    fmt.Println("Leaf:", chunk.MerkleLeafHash)

    fmt.Println()
}

//parallel chunking

	fmt.Println("-- Parallel Chunking --")

	parallelChunks1, err := chunker.ChunklargeFile("test.pdf")
	if err != nil {
		panic(err)
	}

	for _, chunk := range parallelChunks1 {

		fmt.Println("Chunk:", chunk.ChunkIndex)

		fmt.Println("Size:", chunk.ChunkSize)

		fmt.Println("Commitment:", chunk.CommitmentHash)

		fmt.Println("Leaf:", chunk.MerkleLeafHash)

		fmt.Println()
	}

    fmt.Println("Running Sequential Chunking...")

	sequentialChunks, err :=
		chunker.ChunkFile("test.pdf")

	if err != nil {
		log.Fatal(err)
	}

    fmt.Println("Running Parallel Chunking...")

	parallelChunks, err :=
		chunker.ChunklargeFile("test.pdf")

	if err != nil {
		log.Fatal(err)
	}


   
	// Compare Results
	

	fmt.Println()
	fmt.Println("-- Comparing Results --")

	if len(sequentialChunks) != len(parallelChunks) {

		fmt.Println("Chunk count mismatch!")

		fmt.Println(
			"Sequential:",
			len(sequentialChunks),
		)

		fmt.Println(
			"Parallel:",
			len(parallelChunks),
		)

		return
	}

	for i := 0; i < len(sequentialChunks); i++ {

		seq := sequentialChunks[i]

		par := parallelChunks[i]

		fmt.Printf(
			"\nChecking Chunk %d\n",
			i,
		)

		// Compare Chunk Size

		if seq.ChunkSize != par.ChunkSize {

			fmt.Println("Chunk size mismatch")

		} else {

			fmt.Println("Chunk size matches")
		}

		// Compare Commitment

		if seq.CommitmentHash != par.CommitmentHash {

			fmt.Println("Commitment mismatch")

		} else {

			fmt.Println("Commitment matches")
		}

		// Compare Leaf Hash

		if seq.MerkleLeafHash != par.MerkleLeafHash {

			fmt.Println("Leaf hash mismatch")

		} else {

			fmt.Println("Leaf hash matches")
		}
	}

	fmt.Println()
	fmt.Println("Comparison Complete")

}