package domain

import (
	"time"
)

type User struct {
	ID         uint `gorm:"primaryKey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Nickname   string
	PrivateKey string
	PublicKey  string
}

func (User) TableUser() string {
	return "users" // Nombre de la tabla en la base de datos
}
