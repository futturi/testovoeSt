package entites

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	Email        string    `gorm:"unique;not null"`
	RefreshToken string    `gorm:"not null"`
}
