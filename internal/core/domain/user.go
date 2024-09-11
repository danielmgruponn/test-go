package domain

import (
	"time"
)

type User struct {
	ID         string `gorm:"primary_key"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Nickname   string
	PrivateKey string
	PublicKey  string
}

func (User) TableUsers() string {
	return "users"
}
