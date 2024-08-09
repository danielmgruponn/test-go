package domain

import "time"

type User struct {
    ID          uint        `json:"id"`
    NickName    string      `json:"nick_name"`
    Password    string      `json:"-"` // El "-" evita que se serialice en JSON
    CreatedAt   time.Time   `json:"created_at"`
    UpdatedAt   time.Time   `json:"updated_at"`
}

func (User) TableUser() string {
    return "users" // Nombre de la tabla en la base de datos
}