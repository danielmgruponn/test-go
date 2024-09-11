package dto

type UserDTO struct {
	ID         string `json:"id"`
	Nickname   string `json:"nickname"`
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

type UserSafeDTO struct {
	ID        string `json:"id"`
	Nickname  string `json:"nickname"`
	PublicKey string `json:"publicKey"`
}
