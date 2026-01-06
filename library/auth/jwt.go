package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ifulqt/coffeeshops-api/config"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/entity"
)

type Jwt interface {
	GenerateToken(data *entity.JwtData) (string, int, error)
	VerifyToken(token string) (*entity.JwtData, error)
}

type Options struct {
	JwtSecretKey string
	JwtIssuer    string
}

// GenerateToken implements [Jwt].
func (o *Options) GenerateToken(data *entity.JwtData) (string, int, error) {
	now := time.Now().Local()
	expiredAt := now.Add(time.Hour * 24)
	data.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(expiredAt)
	data.RegisteredClaims.Issuer = o.JwtIssuer
	data.RegisteredClaims.NotBefore = jwt.NewNumericDate(now)
	accToken := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	token, err := accToken.SignedString([]byte(o.JwtSecretKey))
	if err != nil {
		return "", 0, err
	}
	return token, int(expiredAt.Unix()), nil
}

// VerifyToken implements [Jwt].
func (o *Options) VerifyToken(token string) (*entity.JwtData, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Signing method invalid")
		}
		return []byte(o.JwtSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if parsedToken.Valid {
		claim, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok || !parsedToken.Valid {
			return nil, err
		}
		jwtData := entity.JwtData{
			UserID: claim["user_id"].(float64),
			Role: claim["role"].(string),
		}
		return &jwtData, nil
	}
	return nil, fmt.Errorf("Token is not invalid")
}

func NewJwt(cfg *config.Config) Jwt {
	opt := new(Options)
	opt.JwtSecretKey = cfg.App.JwtSecretKey
	opt.JwtIssuer = cfg.App.JwtIssuer

	return opt
}
