package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ifulqt/coffeeshops-api/config"
	"github.com/ifulqt/coffeeshops-api/internal/adapter/repository"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/domerror"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/entity"
	"github.com/ifulqt/coffeeshops-api/library/auth"
	"github.com/ifulqt/coffeeshops-api/library/helper"
)

type AuthService interface {
	GetUserByEmail(ctx context.Context, req entity.LoginReq) (*entity.AccessToken, error)
}

type authService struct {
	AuthRepository repository.AuthRepository
	Cfg            *config.Config
	JwtToken       auth.Jwt
}

// GetUserByEmail implements [AuthService].
func (a *authService) GetUserByEmail(ctx context.Context, req entity.LoginReq) (*entity.AccessToken, error) {
	result, err := a.AuthRepository.GetUserByEmail(ctx, req)
	if err != nil {
		code := "[SERVICE] GetUserByEmail - 1"
		log.Errorw(code, err)
		return nil, fmt.Errorf("Get user failed : %w", domerror.ErrUserNotFound)
	}

	err = helper.CheckHashPassword(req.Password, result.Password)
	if err != nil {
		code := "[SERVICE] GetUserByEmail - 2"
		log.Errorw(code, err)
		return nil, fmt.Errorf("Check Hash is failed : %w", domerror.ErrInvalidPassword)
	}

	jwtData := entity.JwtData{
		UserID: float64(result.ID),
		Role:   result.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
			ID:        strconv.Itoa(int(result.ID)),
		},
	}

	accessToken, expiredAt, err := a.JwtToken.GenerateToken(&jwtData)
	if err != nil {
		code := "[SERVICE] GetUserByEmail - 3"
		log.Errorw(code, err)
		return nil, fmt.Errorf("Generate token failed : %w", domerror.ErrGenerateToken)
	}

	resp := entity.AccessToken{
		Token:     accessToken,
		ExpiredAt: int64(expiredAt),
	}

	return &resp, nil
}

func NewAuthService(authRepo repository.AuthRepository, cfg *config.Config, jwtToken auth.Jwt) AuthService {
	return &authService{
		AuthRepository: authRepo,
		Cfg:            cfg,
		JwtToken:       jwtToken,
	}
}
