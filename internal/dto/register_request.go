package dto

type RegisterRequest struct {
	Nickname   string `json:"nickname"`
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}
