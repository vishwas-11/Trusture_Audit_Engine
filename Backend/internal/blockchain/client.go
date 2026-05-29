package blockchain

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

func NewClient(rpcURL string) *ethclient.Client {
	client, err := ethclient.DialContext(context.Background(), rpcURL)
	if err != nil {
		log.Fatal("❌ Failed to connect to blockchain:", err)
	}
	log.Println("✅ Connected to Polygon RPC")
	return client
}
