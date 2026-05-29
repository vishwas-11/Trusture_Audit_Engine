package auth

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func VerifySignature(wallet, nonce, signature string) (bool, error) {

	// Recreate the exact message that was signed
	msg := fmt.Sprintf(
		"Sign this message to login to TRUSTURE:\nNonce: %s\nWallet: %s",
		nonce,
		wallet,
	)

	// Ethereum-specific message prefixing
	hash := crypto.Keccak256Hash([]byte("\x19Ethereum Signed Message:\n" + fmt.Sprint(len(msg)) + msg))

	sigBytes := common.FromHex(signature)

	pubKey, err := crypto.SigToPub(hash.Bytes(), sigBytes)
	if err != nil {
		return false, err
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)

	return recoveredAddr.Hex() == wallet, nil
}
