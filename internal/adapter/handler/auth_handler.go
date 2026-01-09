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

type AuthHandler interface {
	Login(c *fiber.Ctx) error
}

type authHandler struct {
	AuthService service.AuthService
}

// Login implements [AuthHandler].
func (a *authHandler) Login(c *fiber.Ctx) error {
	var err error
	var req request.LoginRequest
	var resp response.AuthResponse
	var errResp response.DefaultErrorResponse

	err = c.BodyParser(&req)
	if err != nil {
		code := "[HANDLER] Login - 1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Invalid request body"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	err = validat.ValidateStruct(&req)
	if err != nil {
		code := "[HANDLER] Login - 2"
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

	reqLogin := entity.LoginReq{
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := a.AuthService.GetUserByEmail(c.Context(), reqLogin)
	if err != nil {
		code := "[HANDLER] Login - 3"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		if errors.Is(err, domerror.ErrUserNotFound) || errors.Is(err, domerror.ErrInvalidPassword) {
			errResp.Meta.Message = "Email or password is incorrect"
			errResp.Meta.Errors = nil
			return c.Status(fiber.StatusUnauthorized).JSON(errResp)
		}
		errResp.Meta.Message = "Internal server error"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	resp.Meta.Status = true
	resp.Meta.Message = "Login successfully"
	resp.AccessToken = result.Token
	resp.ExpiredAt = result.ExpiredAt

	return c.JSON(resp)
}

func NewAuthHandler(authServ service.AuthService) AuthHandler {
	return &authHandler{
		AuthService: authServ,
	}
}
