package dto

type LoginResponse struct {
	ID         string `json:"id"`
	Nickname   string `json:"nickname"`
	Token      string `json:"token"`
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}
