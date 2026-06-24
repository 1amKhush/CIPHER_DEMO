package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"

	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/sha3"
)

var key = []byte("12345678901234567890123456789012")

func main() {

	conn, err := net.Dial("tcp", "localhost:9001")
	if err != nil {
		panic(err)
	}

	message := "Hello World"

	// HASHING
	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(message))
	hashBytes := hash.Sum(nil)

	fmt.Println("Original Hash:", hex.EncodeToString(hashBytes))

	// ENCRYPTION
	aead,err := chacha20poly1305.New(key)
	if err!=nil{
		log.Fatal(err)
	}

	nonce := make([]byte, chacha20poly1305.NonceSize)
	io.ReadFull(rand.Reader, nonce)

	ciphertext := aead.Seal(nil, nonce, []byte(message), nil)

	// SEND
	conn.Write(hashBytes)
	conn.Write(nonce)
	conn.Write(ciphertext)

	fmt.Println("Message Sent")
}