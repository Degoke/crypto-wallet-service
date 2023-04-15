package wallet

import (
	"time"

	"github.com/Degoke/crypto-wallet-service/common"
    "github.com/Degoke/crypto-wallet-service/address"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Wallet struct {
    ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    UserID    uuid.UUID `gorm:"type:uuid;not null"`
    Addresses []address.Address
    Seed     string `gorm:"unique;not null"`
    PrivateKey string `gorm:"unique;not null"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
    Network string `gorm:"not null"`
}

func AutoMigrate(db *gorm.DB) {
    db.AutoMigrate(&Wallet{})
}

func Save(value interface{}) error {
    db := common.GetDB()
    err := db.Save(value).Error
    return err
}

func FindOne(condition interface{}) (Wallet, error) {
    db := common.GetDB()
    var wallet Wallet
    err := db.Where(condition).Preload("Addresses").First(&wallet).Error
    return wallet, err
}

func FindAll(condition interface{}) ([]Wallet, error) {
    db := common.GetDB()
    var wallets []Wallet
    err := db.Where(condition).Preload("Addresses").Find(&wallets).Error
    return wallets, err
}
