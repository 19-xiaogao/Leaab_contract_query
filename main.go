package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"strings"
	"watch_contract/leaab"
)

type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}

func main() {
	const url = "wss://bsc.getblock.io/mainnet/?api_key=5583983f-0c9e-4e6e-b49b-d14523044e3a" // url string
	contractAddress := common.HexToAddress("0xca07f2cADb981c7886a83357B4540002c1F41020")

	rpcClient, err := ethclient.Dial(url)

	if err != nil {
		panic(err)
	}
	query := ethereum.FilterQuery{
		Addresses: []common.Address{
			contractAddress,
		},
	}
	logs := make(chan types.Log)
	sub, err := rpcClient.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}
	contractAbi, err := abi.JSON((strings.NewReader(leeab.LeaabABI)))

	if err != nil {
		log.Fatal(err)
	}

	logTransferSigHash := crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			switch vLog.Topics[0].Hex() {
			case logTransferSigHash.Hex():
				//var transferEvent LogTransfer
				data, err := contractAbi.Unpack("Transfer", vLog.Data)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("------------------Transfer-start-------------------")
				fmt.Println("form", vLog.Topics[1].Hex())
				fmt.Println("to", vLog.Topics[2].Hex())
				fmt.Println("tx value", data)
				fmt.Println("tx hash", vLog.TxHash.Hex())
				fmt.Println("------------------Transfer-end-------------------")
			}
			//fmt.Println("--------------------------------------------------------------------------------------------------------------------------")
			//fmt.Println("contract_address", vLog.Address)
			//fmt.Println("contract_Topics", vLog.Topics)
			//fmt.Println("contract_Data", vLog.Data)
			//fmt.Println("tx hash", vLog.TxHash.Hex())
			//fmt.Println("--------------------------------------------------------------------------------------------------------------------------")
			//fmt.Println(vLog) // pointer to event log
		}
	}

}
