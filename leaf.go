package main

import (
	"encoding/binary"

	"golang.org/x/crypto/sha3"
)

func GenerateLeaf(
	fileID []byte,
	index int,
	chunk []byte,
) []byte {

	indexBytes := make([]byte, 8)
	lengthBytes := make([]byte, 8)

	binary.BigEndian.PutUint64(
		indexBytes,
		uint64(index),
	)

	binary.BigEndian.PutUint64(
		lengthBytes,
		uint64(len(chunk)),
	)

	h := sha3.NewLegacyKeccak256()

	h.Write(fileID)
	h.Write(indexBytes)
	h.Write(lengthBytes)
	h.Write(chunk)
    //leaf is hash of fileid,index,length and data in that chunk
	return h.Sum(nil)
}