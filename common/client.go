package common

import (
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var testnetRpcUrl = "https://data-seed-prebsc-2-s1.binance.org:8545/"

var Client *ethclient.Client

func InitClient() *ethclient.Client {
	client, err := ethclient.Dial(testnetRpcUrl)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	Client = client
	return Client
}

func GetClient() *ethclient.Client {
	return Client
}

func GetContractAddress(currency string) common.Address {
	return common.HexToAddress("0x337610d27c682E347C9cD60BD4b3b107C9d34dDd")
}

func GetTokenABI() string {
	return `[{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"}]`
}