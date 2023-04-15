package transaction

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/Degoke/crypto-wallet-service/address"
	"github.com/Degoke/crypto-wallet-service/common"
	"github.com/Degoke/crypto-wallet-service/wallet"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ec "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
)


func RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/transaction", CreateTransaction)
	// router.GET("/transaction/:id", GetTransaction)
	// router.GET("/transaction", GetTransactions)
	// router.PUT("/transaction/:id", UpdateTransaction)
	// router.DELETE("/transaction/:id", DeleteTransaction)
}

func CreateTransaction(c *gin.Context) {
	validator := NewTransactionValidator()
	if err := validator.Bind(c); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationError(err))
		return
	}

	transactionModel := validator.transactionModel
	fromAddress := ec.HexToAddress(transactionModel.From)
	toAddress := ec.HexToAddress(transactionModel.To)
	value := common.StringTowei(transactionModel.Value)

	client := common.GetClient()

	address, err := address.FindOne(&address.Address{Address: transactionModel.From})
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewError("transaction", err))
		return
	}

	wallet, err := wallet.FindOne(&wallet.Wallet{ID: address.WalletID})
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewError("transaction", err))
		return
	}

	privateKeyECDSA, err := crypto.HexToECDSA(wallet.PrivateKey)

	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewError("transaction", err))
		return
	}

	var data []byte
	var gasLimit uint64 = 21000
	contractAddress := common.GetContractAddress(transactionModel.Currency)

	if transactionModel.Currency != "BNB" {
		tokenAbi, err := abi.JSON(strings.NewReader(common.GetTokenABI()))
		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewError("transaction", err))
			return
		}

		abiData, err := tokenAbi.Pack("transfer", toAddress, value)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewError("transaction", err))
			return
		}

		data = abiData

		gasL, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
			From: fromAddress,
			To:   &contractAddress,
			Data: data,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewError("transaction", err))
			return
		}

		gasLimit = gasL
	}

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewError("transaction", err))
		return
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewError("transaction", err))
		return
	}

	var tx *types.Transaction

	if transactionModel.Currency == "BNB" {
		tx = types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
	} else {
		tx = types.NewTransaction(nonce, contractAddress, value, gasLimit, gasPrice, data)
	}

	

	chainID := big.NewInt(97)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKeyECDSA)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewError("transaction", err))
		return
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewError("transaction", err))
		return
	}

	fmt.Println(signedTx.Hash().String())
	fmt.Println(common.WeiToString(signedTx.GasPrice()))

}