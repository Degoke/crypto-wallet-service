package common

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Seed struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name string `gorm:"unique;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type SeedObject struct {
	Name string 
	Run func(*gorm.DB) error
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&Seed{})
}

func RunSeed(db *gorm.DB, s SeedObject) {
	var seed Seed
	if err := db.Where("name = ?", s.Name).First(&seed).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic("Failed to read seed from database: " + err.Error())
	}

	if err := s.Run(db); err != nil {
		panic("Failed to run seed: " + err.Error())
	}

	seed = Seed{
		Name: seed.Name,
	}

	if err := db.Save(&seed).Error; err != nil {
		panic("Failed to save seed to database: " + err.Error())
	}

}

func RunSeeds(db *gorm.DB, seeds []SeedObject) {
	for _, s := range seeds {
		RunSeed(db, s)
	}
}