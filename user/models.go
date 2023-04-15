package user

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/Degoke/crypto-wallet-service/common"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email     string    `gorm:"not null;unique" json:"email"`
	Password  string    `gorm:"not null" json:"password"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&User{})
}


func (u *User) BeforeSave(tx *gorm.DB) error {
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		fmt.Println("Error:", err)
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword(salt, []byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) ComparePassword(password string) error {
	bytePassword := []byte(password)
	hashedPassword := []byte(u.Password)
	return bcrypt.CompareHashAndPassword(hashedPassword, bytePassword)
}

func Save(value interface{}) error {
	db := common.GetDB()
	err := db.Save(value).Error
	fmt.Println(err)
	return err
}

func (u *User) Delete() error {
	db := common.GetDB()
	err := db.Delete(u).Error
	return err
}

func FindAll() ([]User, error) {
	db := common.GetDB()
	var users []User
	err := db.Find(&users).Error
	return users, err
}

func FindOne(condition interface{}) (User, error) {
	db := common.GetDB()
	var user User
	err := db.Where(condition).First(&user).Error
	return user, err
}
