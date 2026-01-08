package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ifulqt/coffeeshops-api/internal/core/service"
)

type CoffeeShopHandler interface {
	CreateCoffeeShop(c *fiber.Ctx) error
	GetCoffeeShops(c *fiber.Ctx) error
	GetCoffeeShopByID(c *fiber.Ctx) error
	UpdateCoffeeShop(c *fiber.Ctx) error
	DeleteCoffeeShop(c *fiber.Ctx) error
}

type coffeeShopHandler struct {
	CoffeeShopService service.CoffeeShopService
}

// CreateCoffeeShop implements [CoffeeShopHandler].
func (f *coffeeShopHandler) CreateCoffeeShop(c *fiber.Ctx) error {
	panic("unimplemented")
}

// DeleteCoffeeShop implements [CoffeeShopHandler].
func (f *coffeeShopHandler) DeleteCoffeeShop(c *fiber.Ctx) error {
	panic("unimplemented")
}

// GetCoffeeShopByID implements [CoffeeShopHandler].
func (f *coffeeShopHandler) GetCoffeeShopByID(c *fiber.Ctx) error {
	panic("unimplemented")
}

// GetCoffeeShops implements [CoffeeShopHandler].
func (f *coffeeShopHandler) GetCoffeeShops(c *fiber.Ctx) error {
	panic("unimplemented")
}

// UpdateCoffeeShop implements [CoffeeShopHandler].
func (f *coffeeShopHandler) UpdateCoffeeShop(c *fiber.Ctx) error {
	panic("unimplemented")
}

func NewCoffeeShopHandler(coffeeShopServ service.CoffeeShopService) CoffeeShopHandler {
	return &coffeeShopHandler{
		CoffeeShopService: coffeeShopServ,
	}
}
