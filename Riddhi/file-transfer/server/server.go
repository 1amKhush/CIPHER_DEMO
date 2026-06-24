package main

import (
	"bytes"
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

	ln, err := net.Listen("tcp", ":9001")
	if err != nil {
		panic(err)
	}

	fmt.Println("Server Listening...")

	conn, err := ln.Accept()
	if err != nil {
		panic(err)
	}

	// RECEIVE HASH
	receivedHash := make([]byte, 32)
	io.ReadFull(conn, receivedHash)

	// RECEIVE NONCE
	nonce := make([]byte, chacha20poly1305.NonceSize)
	io.ReadFull(conn, nonce)

	// RECEIVE CIPHERTEXT
	buffer := make([]byte, 1024)
	n, _ := conn.Read(buffer)

	ciphertext := buffer[:n]

	// DECRYPT
	aead,err := chacha20poly1305.New(key)
	if err!=nil{
		log.Fatal(err)
	}

	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Received Message:", string(plaintext))

	// GENERATE HASH AGAIN
	hash := sha3.NewLegacyKeccak256()
	hash.Write(plaintext)

	newHash := hash.Sum(nil)

	fmt.Println("Received Hash:", hex.EncodeToString(receivedHash))
	fmt.Println("Generated Hash:", hex.EncodeToString(newHash))

	// VERIFY
	if bytes.Equal(receivedHash, newHash) {
		fmt.Println("AUTHENTIC")
	} else {
		fmt.Println("NOT AUTHENTIC")
	}
}