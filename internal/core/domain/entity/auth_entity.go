package entity

type LoginReq struct {
	Email string
	Password string
}

type AccessToken struct {
	Token string
	ExpiredAt int64
}