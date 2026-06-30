package controllers

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"
)

type Transfer struct {
	Filename  string
	FilePath  string
	ExpiresAt time.Time
}

var (
	store = make(map[string]Transfer)
	mu    sync.Mutex
)

// GeneratePIN creates a random 6-digit PIN and stores the transfer details with a 5-minute expiry.
func GeneratePIN(filename, filePath string) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	var pin string
	for i := 0; i < 10; i++ { // Try up to 10 times to avoid collision
		n, err := rand.Int(rand.Reader, big.NewInt(1000000))
		if err != nil {
			return "", err
		}
		pin = fmt.Sprintf("%06d", n.Int64())
		if _, exists := store[pin]; !exists {
			break
		}
	}

	store[pin] = Transfer{
		Filename:  filename,
		FilePath:  filePath,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	return pin, nil
}

// GetTransfer validates the PIN and returns the transfer details. It returns false if invalid or expired.
// Note: It deletes the PIN upon successful retrieval to ensure single-use.
func GetTransfer(pin string) (Transfer, bool) {
	mu.Lock()
	defer mu.Unlock()

	transfer, exists := store[pin]
	if !exists {
		return Transfer{}, false
	}

	// Check expiration
	if time.Now().After(transfer.ExpiresAt) {
		delete(store, pin)
		return Transfer{}, false
	}

	// Remove PIN from store (one-time use)
	delete(store, pin)
	return transfer, true
}
