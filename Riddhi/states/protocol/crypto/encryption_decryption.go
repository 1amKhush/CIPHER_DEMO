package crypto

import (
	"crypto/rand"
	"io"
	"fmt"

	"golang.org/x/crypto/chacha20poly1305"
)
//Encryption
func Encrypt(
	plaintext []byte,
	key [32]byte,
) (
	nonce [24]byte,
	ciphertext []byte,
	err error,
) {

	aead, err :=
		chacha20poly1305.NewX( //NewX() uses 24-byte nonce
			key[:],
		)

	if err != nil {
		return
	}

	_, err =
		io.ReadFull(
			rand.Reader,
			nonce[:],
		)

	if err != nil {
		return
	}

	fmt.Printf("provider nonce: %x\n", nonce)

	ciphertext =
		aead.Seal(
			nil,
			nonce[:],
			plaintext,
			nil,
		)

	return
}

//Decryption
func Decrypt(
	key [32]byte,
	nonce [24]byte,
	ciphertext []byte,
) ([]byte, error) {

	aead, err :=
		chacha20poly1305.NewX(
			key[:],
		)

	if err != nil {
		return nil, err
	}

	return aead.Open(
		nil,
		nonce[:],
		ciphertext,
		nil,
	)
}



