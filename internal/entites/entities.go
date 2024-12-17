package entites

import (
	"time"
)

type User struct {
	ID        string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Email     string `gorm:"type:varchar(255);uniqueIndex"`
}

type RefreshToken struct {
	ID       uint   `gorm:"primaryKey"`
	UserID   string `gorm:"type:uuid;index"`
	Hash     string `gorm:"type:text"`
	IssuedAt time.Time
	ClientIp string `gorm:"type:varchar(255)"`
	Used     bool
}
