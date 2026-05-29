package blockchain

import "math/big"

type DonationReceivedEvent struct {
	Donor     string
	NGO       string
	Amount    *big.Int
	Purpose   string
	ProofHash [32]byte
	Timestamp *big.Int
}
