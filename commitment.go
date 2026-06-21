package main

import (
	"crypto/rand"

	"golang.org/x/crypto/sha3"
)

func GenerateCommitment(chunk []byte) ([]byte, []byte, error) {

	key := make([]byte, 32)

	_, err := rand.Read(key)
	if err != nil {
		return nil, nil, err
	}

	h := sha3.NewLegacyKeccak256()

	h.Write(key)
	h.Write(chunk)

	commitment := h.Sum(nil)
    // commitment is keys||data in chunk
	return key, commitment, nil
}