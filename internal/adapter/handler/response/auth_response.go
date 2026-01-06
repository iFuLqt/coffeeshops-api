package response

type AuthResponse struct {
	Meta Meta `json:"meta"`
	AccessToken string `json:"access_token"`
	ExpiredAt int64 `json:"expired_at"`
}