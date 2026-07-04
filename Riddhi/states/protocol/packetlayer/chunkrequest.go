package packetlayer

type FileMetadataRequestPayload struct {
	FileName string
	//FileID	[32]byte //for future
}

type ChunkRequestPayload struct {
	ChunkIndex uint32
	Nonce      [16]byte
}