package dto

type RegisterRequest struct {
	NickName   string `json:"nickName"`
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}
