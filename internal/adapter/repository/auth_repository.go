package repository

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/entity"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/model"
	"gorm.io/gorm"
)

type AuthRepository interface {
	GetUserByEmail(ctx context.Context, req entity.LoginReq) (*entity.UserEntity, error)
}

type authRepository struct {
	db *gorm.DB
}

// GetUserByEmail implements [AuthRepository].
func (a *authRepository) GetUserByEmail(ctx context.Context, req entity.LoginReq) (*entity.UserEntity, error) {
	var modelUser model.User

	err := a.db.WithContext(ctx).Where("email = ?", req.Email).First(&modelUser).Error
	if err != nil {
		code := "[REPOSITORY] GetUserByEmail - 1"
		log.Errorw(code, err)
		return nil, err
	}

	userEntity := entity.UserEntity{
		ID: modelUser.ID,
		Name: modelUser.Name,
		Email: modelUser.Email,
		Password: modelUser.Password,
		Role: modelUser.Role,
	}

	return &userEntity, nil
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{
		db: db,
	}
}
