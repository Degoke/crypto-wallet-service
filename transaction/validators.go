package transaction

import (

	"github.com/Degoke/crypto-wallet-service/common"
	"github.com/gin-gonic/gin"
)

type TransactionValidator struct {
	Transaction struct {
		From string `json:"from" binding:"required"`
		To   string `json:"to" binding:"required"`
		Value string `json:"value" binding:"required"`
		Currency string `json:"currency" binding:"required"`
	} `json:"transaction"`
	transactionModel Transaction `json:"-"`
}

func (t *TransactionValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, t)
	if err != nil {
		return err
	}

	t.transactionModel.From = t.Transaction.From
	t.transactionModel.To = t.Transaction.To
	t.transactionModel.Value = t.Transaction.Value
	t.transactionModel.Currency = t.Transaction.Currency

	return nil
}

func NewTransactionValidator() TransactionValidator {
	transactionValidator := TransactionValidator{}
	return transactionValidator
}