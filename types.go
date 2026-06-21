package main
//data of chunk that we are collecting and calculating
type ChunkInfo struct {
	Index      int
	Size       int
	Key        []byte
	Data       []byte
	Commitment []byte
	Leaf       []byte
}