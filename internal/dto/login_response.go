package dto

type LoginResponse struct {
	ID         uint   `json:"id"`
	Nickname   string `json:"nickname"`
	Token      string `json:"token"`
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}
