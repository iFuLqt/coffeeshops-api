package service

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/ifulqt/coffeeshops-api/internal/adapter/repository"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/domerror"
	"github.com/ifulqt/coffeeshops-api/library/helper"
)

type UserService interface {
	UpdatePassword(ctx context.Context, newPass string, currentPass string, id int64) error
}

type userService struct {
	UserRepository repository.UserRepository
}

// UpdatePassword implements [UserService].
func (u *userService) UpdatePassword(ctx context.Context, newPass string, currentPass string, id int64) error {
	result, err := u.UserRepository.GetUserByID(ctx, id)
	if err != nil {
		code := "[SERVICE] UpdatePassword - 1"
		log.Errorw(code, err)
		return domerror.ErrUserNotFound
	}

	if newPass == currentPass {
		return domerror.ErrSamePassword
	}

	err = helper.CheckHashPassword(currentPass, result.Password)
	if err != nil {
		code := "[SERVICE] UpdatePassword - 2"
		log.Errorw(code, err)
		return domerror.ErrInvalidPassword
	}

	hashNewPass, err := helper.HashPassword(newPass)
	if err != nil {
		code := "[SERVICE] UpdatePassword - 3"
		log.Errorw(code, err)
		return err
	}

	return u.UserRepository.UpdatePassword(ctx, hashNewPass, id)
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		UserRepository: userRepo,
	}
}
