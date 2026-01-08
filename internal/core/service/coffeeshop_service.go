package service

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/ifulqt/coffeeshops-api/internal/adapter/repository"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/entity"
)

type CoffeeShopService interface {
	CreateCoffeeShop(ctx context.Context, req entity.CoffeeShopEntity) error
	GetCoffeeShops(ctx context.Context) ([]entity.CoffeeShopEntity, error)
	GetCoffeeShopByID(ctx context.Context, id int) (*entity.CoffeeShopEntity, error)
	UpdateCoffeeShop(ctx context.Context, req entity.CoffeeShopEntity) error
	DeleteCoffeeShop(ctx context.Context, id int) error
}

type coffeeShopService struct {
	CoffeeShopRepository repository.CoffeeShopRepository
}

// CreateCoffeeShop implements [CoffeeShopService].
func (c *coffeeShopService) CreateCoffeeShop(ctx context.Context, req entity.CoffeeShopEntity) error {
	err := c.CoffeeShopRepository.CreateCoffeeShop(ctx, req)
	if err != nil {
		code := "[HANDLER] CreateCoffeeShop - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// DeleteCoffeeShop implements [CoffeeShopService].
func (c *coffeeShopService) DeleteCoffeeShop(ctx context.Context, id int) error {
	panic("unimplemented")
}

// GetCoffeeShopByID implements [CoffeeShopService].
func (c *coffeeShopService) GetCoffeeShopByID(ctx context.Context, id int) (*entity.CoffeeShopEntity, error) {
	panic("unimplemented")
}

// GetCoffeeShops implements [CoffeeShopService].
func (c *coffeeShopService) GetCoffeeShops(ctx context.Context) ([]entity.CoffeeShopEntity, error) {
	panic("unimplemented")
}

// UpdateCoffeeShop implements [CoffeeShopService].
func (c *coffeeShopService) UpdateCoffeeShop(ctx context.Context, req entity.CoffeeShopEntity) error {
	panic("unimplemented")
}

func NewCoffeeShopService(coffeeShopRepo repository.CoffeeShopRepository) CoffeeShopService {
	return &coffeeShopService{
		CoffeeShopRepository: coffeeShopRepo,
	}
}
