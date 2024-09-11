package dto

type RegisterRequest struct {
	ID         string `json:"id"`
	Nickname   string `json:"nickname"`
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}
