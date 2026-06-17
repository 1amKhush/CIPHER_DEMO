package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"io"

	"golang.org/x/crypto/sha3"
	"github.com/ethereum/go-ethereum/crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/chacha20poly1305"

)

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}


func main (){

data, err := os.ReadFile("file.txt") // For read access.
if err != nil {
	log.Fatal(err)
}
fmt.Println(data)
fmt.Println(data[0:2]) //0:2-this would print data[0],data[1] not data[2]
fmt.Println(data[5:8],data[9:13])
fmt.Println(string(data[0:2]))
fmt.Println(string(data[5:9]),string(data[9:13]))
fmt.Println(string(data))

hash:=sha3.NewLegacyKeccak256()
hash.Write(data)
result:=hash.Sum(nil)
fmt.Printf("%x\n", result) //readable form(hexadecimal number)
fmt.Println(result) //in bytes
fmt.Println(string(result))
fmt.Println(hex.EncodeToString(result))//readable form(hexadecimal number)

//other way of keccka use

fmt.Print(crypto.Keccak256(data))
fmt.Printf("\n")
fmt.Println(hex.EncodeToString(crypto.Keccak256(data)))
fmt.Print(crypto.Keccak256Hash(data))
fmt.Printf("\n")

//encryption
// Load your secret key from a safe place and reuse it across multiple
	// Seal/Open calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	// When decoded the key should be 16 bytes (AES-128) or 32 (AES-256).
	key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
	fmt.Printf("key:%x",key)
	fmt.Printf("\n")
	plaintext := []byte("exampleplaintext")

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	fmt.Printf("nonce:%x\n", nonce)

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	fmt.Printf("Encrypted:%x\n", ciphertext)


    decrypted, err := aesgcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        panic(err)
    }

    fmt.Println("Decrypted:", string(decrypted))


	//otherway-chacha20
	keynew := make([]byte, chacha20poly1305.KeySize)

	if _, err := io.ReadFull(rand.Reader, keynew); err != nil {
		panic(err)
	}

	aead, err := chacha20poly1305.New(keynew)
	if err != nil {
		panic(err)
	}

	noncenew := make([]byte, chacha20poly1305.NonceSize)

	_, err = io.ReadFull(rand.Reader, noncenew)
	if err != nil {
		panic(err)
	}

	plaintextnew := []byte("exampleplaintext")

	ciphertextnew := aead.Seal(nil, noncenew, plaintextnew, nil)

	fmt.Printf("Ciphertext: %x\n", ciphertextnew)

	decryptednew, err := aead.Open(nil, noncenew, ciphertextnew, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Plaintext:", string(decryptednew))

//line to line reading
readFile, err := os.Open("file.txt")
	if err != nil {
		panic(err)
	}
	defer readFile.Close()
scanner:=bufio.NewScanner(readFile)
for scanner.Scan() {
		fmt.Println("Scanned Line:", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}


}