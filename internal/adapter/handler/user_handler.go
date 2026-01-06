package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/ifulqt/coffeeshops-api/internal/adapter/handler/request"
	"github.com/ifulqt/coffeeshops-api/internal/adapter/handler/response"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/domerror"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/entity"
	"github.com/ifulqt/coffeeshops-api/internal/core/service"
	"github.com/ifulqt/coffeeshops-api/library/validat"
)

type UserHandler interface {
	UpdatePassword(c *fiber.Ctx) error
}

type userHandler struct {
	UserService service.UserService
}

// UpdatePassword implements [UserHandler].
func (u *userHandler) UpdatePassword(c *fiber.Ctx) error {
	var err error
	var errResp response.DefaultErrorResponse
	var resp response.DefaultSuccessResponse
	var req request.UpdatePasswordRequest

	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		errResp.Meta.Status = false
		errResp.Meta.Message = "Unauthorized Access"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusUnauthorized).JSON(errResp)
	}

	err = c.BodyParser(&req)
	if err != nil {
		code := "[HANDLER] UpdatePassword - 1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Invalid request body"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	err = validat.ValidateStruct(&req)
	if err != nil {
		code := "[HANDLER] UpdatePassword - 2"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Invalid request data"
		val, ok := err.(validat.ValidationError)
		if ok {
			errResp.Meta.Errors = val.Message
		} else {
			errResp.Meta.Errors = nil
		}
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	err = u.UserService.UpdatePassword(c.Context(), req.NewPassword, req.CurrentPassword, int(userID))
	if err != nil {
		code := "[HANDLER] UpdatePassword - 3"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		if errors.Is(err, domerror.ErrInvalidPassword) || errors.Is(err, domerror.ErrUserNotFound) {
			errResp.Meta.Message = "Current password is wrong"
			errResp.Meta.Errors = nil
			return c.Status(fiber.StatusBadRequest).JSON(errResp)
		} else if errors.Is(err, domerror.ErrSamePassword) {
			errResp.Meta.Message = "New password mush be different from current password"
			errResp.Meta.Errors = nil
			return c.Status(fiber.StatusBadRequest).JSON(errResp)
		}
		errResp.Meta.Message = "Internal server error"
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)

	}

	resp.Meta.Status = true
	resp.Meta.Message = "Update Password Success"
	resp.Meta.Errors = nil

	return c.JSON(resp)
}

func NewUserHandler(userServ service.UserService) UserHandler {
	return &userHandler{
		UserService: userServ,
	}
}
