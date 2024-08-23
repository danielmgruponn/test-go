package dto

type UserDTO struct {
	ID         uint   `json:"id"`
	NickName   string `json:"nickName"`
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}
