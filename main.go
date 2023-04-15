package main

import (
	"log"

	"github.com/Degoke/crypto-wallet-service/address"
	"github.com/Degoke/crypto-wallet-service/common"
	"github.com/Degoke/crypto-wallet-service/transaction"
	"github.com/Degoke/crypto-wallet-service/user"
	"github.com/Degoke/crypto-wallet-service/wallet"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	user.AutoMigrate(db)
	wallet.AutoMigrate(db)
	address.AutoMigrate(db)
}

func main() {
	common.LoadENV()
	db := common.InitDB()
	Migrate(db)

	postgresdb, err := db.DB()
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
		panic(err)
	}

	defer postgresdb.Close()

	client := common.InitClient()
	defer client.Close()

	r := gin.Default()
	v1 := r.Group("/api/v1")
	user.RegisterUnprotectedRoutes(v1.Group("/users"))

	v1.Use(user.AuthMiddleware(true))

	user.RegisterRoutes(v1.Group("/users"))
	wallet.RegisterRoutes(v1.Group("/wallets"))
	transaction.RegisterRoutes(v1.Group("/transactions"))

	port := common.GetEnv("PORT")
	r.Run(":" + port)
}