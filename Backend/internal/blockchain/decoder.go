package blockchain

import (
	// "log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func DecodeDonationEvent(vLog types.Log) (*DonationReceivedEvent, error) {
	abiBytes, err := os.ReadFile("internal/blockchain/contracts/trusture_abi.json")
	if err != nil {
		return nil, err
	}

	contractABI, err := abi.JSON(strings.NewReader(string(abiBytes)))
	if err != nil {
		return nil, err
	}

	event := new(DonationReceivedEvent)

	err = contractABI.UnpackIntoInterface(event, "DonationReceived", vLog.Data)
	if err != nil {
		return nil, err
	}

	// Indexed fields
	event.Donor = common.HexToAddress(vLog.Topics[1].Hex()).Hex()
	event.NGO = common.HexToAddress(vLog.Topics[2].Hex()).Hex()

	return event, nil
}
