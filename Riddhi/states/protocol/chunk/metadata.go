package chunk

type ChunkMeta struct {
    ChunkIndex int
    ChunkSize int
    ChunkData []byte
    ChunkHash [32]byte
    LeafHash [32]byte
}

type DefaultFileChunker struct {
 ChunkSize int
}

type ChunkJob struct {
    Index  int
    Offset int64
}

type ChunkedFile struct { 
    FileID [32]byte 

    FileSize int64
	ChunkSize int
    
    Chunks []ChunkMeta
     LeafHashes [][32]byte 
     MerkleRoot [32]byte 
    }
