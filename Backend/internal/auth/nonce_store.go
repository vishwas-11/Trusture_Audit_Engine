package auth

import (
	"sync"
	"time"
)

type nonceEntry struct {
	Nonce     string
	ExpiresAt time.Time
}

var (
	nonces = make(map[string]nonceEntry)
	mu     sync.Mutex
)

func SetNonce(wallet, nonce string) {
	mu.Lock()
	defer mu.Unlock()
	
	nonces[wallet] = nonceEntry{
		Nonce:     nonce,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
}

func GetNonce(wallet string) (string, bool) {
	mu.Lock()
	defer mu.Unlock()

	entry, exists := nonces[wallet]
	if !exists || time.Now().After(entry.ExpiresAt) {
		return "", false
	}

	return entry.Nonce, true
}

func DeleteNonce(wallet string) {
	mu.Lock()
	defer mu.Unlock()
	delete(nonces, wallet)
}
