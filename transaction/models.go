package transaction

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    From  string    `gorm:"not null"`
	To  string    `gorm:"not null"`
	Currency  string    `gorm:"not null"`
	Value string   `gorm:"not null"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}