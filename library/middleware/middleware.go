package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ifulqt/coffeeshops-api/internal/adapter/handler/response"
	"github.com/ifulqt/coffeeshops-api/library/auth"
)

type Middleware interface {
	CheckToken() fiber.Handler
}

type Options struct {
	AuthJwt auth.Jwt
}

// CheckToken implements [Middleware].
func (o *Options) CheckToken() func(c *fiber.Ctx) error {
	return func (c *fiber.Ctx) error {
		var errResp response.DefaultErrorResponse
		authHandler := c.Get("Authorization")
		if authHandler == "" {
			errResp.Meta.Status = false
			errResp.Meta.Message = "Missing authorization header"
			return c.Status(fiber.StatusUnauthorized).JSON(errResp)
		}
		
		tokenString := strings.Split(authHandler, "Bearer ")[1]
		claims, err := o.AuthJwt.VerifyToken(tokenString)
		if err != nil {
			errResp.Meta.Status = false
			errResp.Meta.Message = "Invalid token"
			return c.Status(fiber.StatusBadRequest).JSON(errResp)
		}
		c.Locals("user", claims)

		return c.Next()
	}
}

func NewMiddleware(authJwt auth.Jwt) Middleware {
	return &Options{
		AuthJwt: authJwt,
	}
}
