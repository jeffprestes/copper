package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jeffprestes/copper/contracts"
)

// LogFill represents a LogFill event raised by the Exchange contract.
type LogFill struct {
	Maker                  common.Address
	Taker                  common.Address
	FeeRecipient           common.Address
	MakerToken             common.Address
	TakerToken             common.Address
	FilledMakerTokenAmount *big.Int
	FilledTakerTokenAmount *big.Int
	PaidMakerFee           *big.Int
	PaidTakerFee           *big.Int
	Tokens                 [32]byte
	OrderHash              [32]byte
	Raw                    types.Log // Blockchain specific contextual infos
}

func main() {
	//Connect to mainnet
	client, err := ethclient.Dial("https://mainnet.infura.io/QPF0qjGpH9OjFuuMrCse")
	if err != nil {
		log.Fatal(err)
	}

	//Input address of 0x v1 contract and build event query
	contractAddress := common.HexToAddress("0x12459C951127e0c374FF9105DdA097662A027093")
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(6338488),
		ToBlock:   big.NewInt(6338504),
		Addresses: []common.Address{
			contractAddress,
		},
	}

	//Run query against ethereum client
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(contracts.ExchangeABI)))
	if err != nil {
		log.Fatal(err)
	}

	for _, vLog := range logs {
		log.Printf("Log: %d\n", vLog.Index)
		fillEvent := LogFill{}
		err := contractAbi.Unpack(&fillEvent, "LogFill", vLog.Data)
		if err != nil {
			fmt.Println(" ")
			log.Println(err)
			fmt.Println(" ")
		}
	}
}
