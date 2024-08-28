package dto

type UserDTO struct {
	ID         uint   `json:"id"`
	Nickname   string `json:"nickname"`
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}

type UserSafeDTO struct {
	ID        uint   `json:"id"`
	Nickname  string `json:"nickname"`
	PublicKey string `json:"publicKey"`
}
