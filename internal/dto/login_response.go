package dto

type LoginResponse struct {
	ID         uint   `json:"id"`
	NickName   string `json:"nickName"`
	Token      string `json:"token"`
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}
