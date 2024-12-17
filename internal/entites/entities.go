package entites

import (
	"time"
)

type User struct {
	ID        string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UpdatedAt time.Time
	Email     string `gorm:"type:varchar(255);uniqueIndex"`
	RefreshToken string `gorm:"type:text"`
	ClientIp string `gorm:"type:varchar(255)"`
}