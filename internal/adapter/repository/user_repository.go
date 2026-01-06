package repository

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/entity"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	UpdatePassword(ctx context.Context, newPass string, id int) error
	GetUserByID(ctx context.Context, id int) (*entity.UserEntity, error)
}

type userRepository struct {
	db *gorm.DB
}

// CheckCurrentPassword implements [UserRepository].
func (u *userRepository) GetUserByID(ctx context.Context, id int) (*entity.UserEntity, error) {
	var modelUser model.User
	err := u.db.Where("id = ?", id).First(&modelUser).Error
	if err != nil {
		code := "[REPOSITORY] GetUserByID - 1"
		log.Errorw(code, err)
		return nil, err
	}

	userEntity := entity.UserEntity{
		ID:       modelUser.ID,
		Name:     modelUser.Name,
		Email:    modelUser.Email,
		Password: modelUser.Password,
		Role:     modelUser.Role,
	}

	return &userEntity, nil
}

// UpdatePassword implements [UserRepository].
func (u *userRepository) UpdatePassword(ctx context.Context, newPass string, id int) error {
	err := u.db.Model(&model.User{}).Where("id = ?", id).Update("password", newPass).Error
	if err != nil {
		code := "[REPOSITORY] UpdatePassword - 2"
		log.Errorw(code, err)
		return err
	}
	return nil
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}
