package blockchain

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func ListenDonations(
	client *ethclient.Client,
	contractAddress string,
	handler func(types.Log),
) {
	address := common.HexToAddress(contractAddress)
	startBlock := int64(33600000)
	seenLogs := make(map[string]struct{})
	lastProcessedBlock := uint64(0)

	baseQuery := ethereum.FilterQuery{
		Addresses: []common.Address{address},
	}

	log.Printf("🎯 Donation listener configured | contract=%s startBlock=%d", address.Hex(), startBlock)

	processIfNew := func(vLog types.Log, source string) {
		key := fmt.Sprintf("%s:%d", vLog.TxHash.Hex(), vLog.Index)
		if _, exists := seenLogs[key]; exists {
			return
		}
		seenLogs[key] = struct{}{}

		if vLog.BlockNumber > lastProcessedBlock {
			lastProcessedBlock = vLog.BlockNumber
		}

		log.Printf("📥 New blockchain event received via %s | tx=%s block=%d", source, vLog.TxHash.Hex(), vLog.BlockNumber)
		handler(vLog)
	}

	// --- STEP 1: FETCH PAST EVENTS (BACKFILL) ---
	backfillQuery := ethereum.FilterQuery{
		Addresses: baseQuery.Addresses,
		FromBlock: big.NewInt(startBlock),
	}

	pastLogs, err := client.FilterLogs(context.Background(), backfillQuery)
	if err != nil {
		log.Println("⚠️ Error fetching past logs:", err)
	} else {
		log.Println("📜 Backfilling past events:", len(pastLogs))
		for _, vLog := range pastLogs {
			processIfNew(vLog, "backfill")
		}
	}

	// If there were no backfilled events, bootstrap from current chain head.
	if lastProcessedBlock == 0 {
		headBlock, headErr := client.BlockNumber(context.Background())
		if headErr != nil {
			log.Println("⚠️ Could not read current chain head:", headErr)
		} else {
			lastProcessedBlock = headBlock
			log.Printf("⛓️ Listener head set to current block=%d", headBlock)
		}
	}

	// --- STEP 2: SUBSCRIBE TO NEW EVENTS (WS) ---
	ch := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), baseQuery, ch)
	if err != nil {
		log.Fatal("❌ Subscription error:", err)
	}

	log.Println("📡 Listening for Donation events (websocket + poll fallback)...")

	// --- STEP 3: POLL FALLBACK (covers dropped WS events) ---
	pollTicker := time.NewTicker(12 * time.Second)
	defer pollTicker.Stop()

	for {
		select {
		case err := <-sub.Err():
			log.Println("❌ Listener websocket error:", err)

		case vLog := <-ch:
			processIfNew(vLog, "websocket")

		case <-pollTicker.C:
			latestBlock, err := client.BlockNumber(context.Background())
			if err != nil {
				log.Println("⚠️ Poll error reading latest block:", err)
				continue
			}

			if latestBlock <= lastProcessedBlock {
				continue
			}

			fromBlock := lastProcessedBlock + 1
			pollQuery := ethereum.FilterQuery{
				Addresses: baseQuery.Addresses,
				FromBlock: new(big.Int).SetUint64(fromBlock),
				ToBlock:   new(big.Int).SetUint64(latestBlock),
			}

			logs, err := client.FilterLogs(context.Background(), pollQuery)
			if err != nil {
				log.Printf("⚠️ Poll logs fetch failed for range %d-%d: %v", fromBlock, latestBlock, err)
				continue
			}

			if len(logs) == 0 {
				lastProcessedBlock = latestBlock
				continue
			}

			for _, vLog := range logs {
				processIfNew(vLog, "poll")
			}
		}
	}
}
