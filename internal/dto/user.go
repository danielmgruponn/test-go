package dto

type UserDTO struct {
	ID         uint   `json:"id"`
	Nickname   string `json:"nickname"`
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

type UserSafeDTO struct {
	ID        uint   `json:"id"`
	Nickname  string `json:"nickname"`
	PublicKey string `json:"publicKey"`
}
