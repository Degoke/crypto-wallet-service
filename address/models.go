package address
import (
	"time"

	"github.com/Degoke/crypto-wallet-service/common"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Address struct {
    ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    Address  string    `gorm:"unique;not null"`
    WalletID uuid.UUID `gorm:"type:uuid;not null"`
    Currency string    `gorm:"not null"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func AutoMigrate(db *gorm.DB) {
    db.AutoMigrate(&Address{})
}

func Save(value interface{}) error {
    db := common.GetDB()
    err := db.Save(value).Error
    return err
}

func FindOne(condition interface{}) (Address, error) {
    db := common.GetDB()
    var address Address
    err := db.Where(condition).First(&address).Error
    return address, err
}

func FindAll(condition interface{}) ([]Address, error) {
    db := common.GetDB()
    var wallets []Address
    err := db.Where(condition).Find(&wallets).Error
    return wallets, err
}