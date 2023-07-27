package user

type UserRegisterResp struct {
	HashId string `json:"hash_id"`
	Token  string `json:"token"`
}

type UserLoginResp struct {
	HashId string `json:"hash_id"`
	Token  string `json:"token"`
}
