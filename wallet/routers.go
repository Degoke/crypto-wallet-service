package wallet

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/Degoke/crypto-wallet-service/address"
	"github.com/Degoke/crypto-wallet-service/common"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	// "github.com/bnb-chain/go-sdk/client"
	// "github.com/bnb-chain/go-sdk/client/rpc"
	"github.com/bnb-chain/go-sdk/keys"
	// "github.com/bnb-chain/go-sdk/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ec "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
)

func RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/", CreateWallet)
	router.GET("/", GetWallets)
	router.POST("/address/:walletId/:currency", createAddress)
	router.GET("/balance/:id", getBalance)
	router.GET("/:id", GetWallet)

}

func CreateWallet(c *gin.Context) {
	userId := c.MustGet("userId").(uuid.UUID)
	fmt.Println(userId)
	seed := generateSeed()

	keyManager, err := keys.NewMnemonicKeyManager(seed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewError("wallet", err))
		return
	}

	mnemonic, _ := keyManager.ExportAsMnemonic()
	privateKey, _ := keyManager.ExportAsPrivateKey()
	mainAddr := keyManager.GetAddr()
	
	walletID := uuid.New()

	wallet := Wallet{
		ID:   walletID,
		UserID:  userId,
		Seed:      mnemonic,
		PrivateKey: privateKey,
		Network: "BSC",
	}

	if err := Save(&wallet); err != nil {
		c.JSON(http.StatusInternalServerError, common.NewError("wallet", err))
		return
	}

	newAddress := address.Address{
		Address: mainAddr.String(),
		WalletID: walletID,
		Currency: "BNB",
	}

	if err := address.Save(&newAddress); err != nil {
		c.JSON(http.StatusInternalServerError, common.NewError("wallet", err))
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"message": "Wallet created"})

}

func GetWallets(c *gin.Context) {
	userId := c.MustGet("userId").(uuid.UUID)
	wallets, err := FindAll(&Wallet{UserID: userId})

	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewError("wallet", err))
		return
	}

	serializer := WalletsSerializer{wallets}

	c.JSON(http.StatusOK, gin.H{"wallets": serializer.Response()})
}

func GetWallet(c *gin.Context) {
	walletIdString := c.Param("id")
	walletId, err := uuid.Parse(walletIdString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid wallet ID",
		})
		return
	}
	userId := c.MustGet("userId").(uuid.UUID)

	wallet, err := FindOne(&Wallet{ID: walletId, UserID: userId})

	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewError("wallet", err))
		return
	}

	serializer := WalletSerializer{wallet}

	c.JSON(http.StatusOK, gin.H{"wallet": serializer.Response()})
}

func generateSeed() string {
	entropy, _ := bip39.NewEntropy(256)
	seed, _ := bip39.NewMnemonic(entropy)

	return seed
}
func createAddress(c *gin.Context) {
	currency := c.Param("currency")
	walletIdString := c.Param("walletId")
	if currency != "BUSD" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Currency not supported",
		})
		return
	}

	walletId, err := uuid.Parse(walletIdString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid wallet ID",
		})
		return
	}
	userId := c.MustGet("userId").(uuid.UUID)

	wallet, err := FindOne(&Wallet{ID: walletId, UserID: userId})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error finding wallet",
		})
		return
	}

	

	addr := generateAddress(wallet.PrivateKey)

	newAddress := address.Address{
		Address: addr.String(),
		WalletID: wallet.ID,
		Currency: currency,
	}

	if err := address.Save(&newAddress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error saving address",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"currency": currency,
		"address": addr.String(),
	})
}

func getBalance(c *gin.Context) {
	walletIdString := c.Param("id")
	walletId, err := uuid.Parse(walletIdString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid address ID",
		})
		return
	}

	userId := c.MustGet("userId").(uuid.UUID)

	wallet, err := FindOne(&Wallet{ID: walletId, UserID: userId})

	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewError("wallet", err))
		return
	}

	balances := make(map[string]*big.Float)


	for _, addr := range wallet.Addresses {
		if addr.Currency == "BNB" {
			balance, err := getBNBbalance(ec.HexToAddress(addr.Address))
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.NewError("wallet", err))
			return
		}
		balances[addr.Currency] = balance
		} else {
			balance, err := getTokenBalance(ec.HexToAddress(addr.Address))
			if err != nil {
				c.JSON(http.StatusInternalServerError, common.NewError("wallet", err))
				return
			}
			balances[addr.Currency] = balance
		}
	}

	c.JSON(http.StatusOK, gin.H{"balances": balances})


}

func generateAddress(privateKey string) ec.Address {
	privateKeyECDSA, _ := crypto.HexToECDSA(privateKey)
	publicKey := crypto.PubkeyToAddress(privateKeyECDSA.PublicKey)
	return publicKey
}

func getBNBbalance(address ec.Address) (*big.Float, error) {
	client := common.GetClient()
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return nil, err
	}
	balanceFloat := new(big.Float).SetInt(balance)
	bnbBalance := new(big.Float).Quo(balanceFloat, big.NewFloat(1e18))
	return bnbBalance, nil
}

func getTokenBalance(address ec.Address) (*big.Float, error) {
	tokenAbi, err := abi.JSON(strings.NewReader(common.GetTokenABI()))
	if err != nil {
		return nil, err
	}

	callData, err := tokenAbi.Pack("balanceOf", address)
	if err != nil {
		return nil, err
	}

	contractAddress := common.GetContractAddress("USDT")

	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: callData,
	}

	client := common.GetClient()

	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}

	var balance *big.Int
	err = tokenAbi.Unpack(&balance, "balanceOf", result)
	if err != nil {
		return nil, err
	}

	balanceFloat := new(big.Float).SetInt(balance)
	tokenBalance := new(big.Float).Quo(balanceFloat, big.NewFloat(1e18))
	return tokenBalance, nil
}