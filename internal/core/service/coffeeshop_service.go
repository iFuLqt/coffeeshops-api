package service

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/ifulqt/coffeeshops-api/internal/adapter/repository"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/entity"
	"github.com/ifulqt/coffeeshops-api/library/helper"
)

type CoffeeShopService interface {
	CreateCoffeeShop(ctx context.Context, req entity.CoffeeShopEntity) (int64, error)
	GetCoffeeShops(ctx context.Context) ([]entity.CoffeeShopEntity, error)
	GetCoffeeShopByID(ctx context.Context, id int64) (*entity.CoffeeShopEntity, error)
	UpdateCoffeeShop(ctx context.Context, req entity.CoffeeShopEntity, idCoffeeShop int64) error
	DeleteCoffeeShop(ctx context.Context, id int64) error
}

type coffeeShopService struct {
	CoffeeShopRepository repository.CoffeeShopRepository
}

// CreateCoffeeShop implements [CoffeeShopService].
func (c *coffeeShopService) CreateCoffeeShop(ctx context.Context, req entity.CoffeeShopEntity) (int64, error) {
	slug := helper.GenerateSlug(req.Name)
	req.Slug = slug

	id, err := c.CoffeeShopRepository.CreateCoffeeShop(ctx, req)
	if err != nil {
		code := "[SERVICE] CreateCoffeeShop - 3"
		log.Errorw(code, err)
		return 0, err
	}
	return id, nil
}

// DeleteCoffeeShop implements [CoffeeShopService].
func (c *coffeeShopService) DeleteCoffeeShop(ctx context.Context, id int64) error {
	err := c.CoffeeShopRepository.DeleteCoffeeShop(ctx, id)
	if err != nil {
		code := "[SERVICE] DeleteCoffeeShop - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// GetCoffeeShopByID implements [CoffeeShopService].
func (c *coffeeShopService) GetCoffeeShopByID(ctx context.Context, id int64) (*entity.CoffeeShopEntity, error) {
	result, err := c.CoffeeShopRepository.GetCoffeeShopByID(ctx, id)
	if err != nil {
		code := "[SERVICE] GetCoffeeShops - 1"
		log.Errorw(code, err)
		return nil, err
	}
	return result, nil
}

// GetCoffeeShops implements [CoffeeShopService].
func (c *coffeeShopService) GetCoffeeShops(ctx context.Context) ([]entity.CoffeeShopEntity, error) {
	results, err := c.CoffeeShopRepository.GetCoffeeShops(ctx)
	if err != nil {
		code := "[SERVICE] GetCoffeeShops - 1"
		log.Errorw(code, err)
		return nil, err
	}
	return results, nil
}

// UpdateCoffeeShop implements [CoffeeShopService].
func (c *coffeeShopService) UpdateCoffeeShop(ctx context.Context, req entity.CoffeeShopEntity, idCoffeeShop int64) error {
	slug := helper.GenerateSlug(req.Name)
	req.Slug = slug

	err := c.CoffeeShopRepository.UpdateCoffeeShop(ctx, req, idCoffeeShop)
	if err != nil {
		code := "[SERVICE] UpdateCoffeeShop - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

func NewCoffeeShopService(coffeeShopRepo repository.CoffeeShopRepository) CoffeeShopService {
	return &coffeeShopService{
		CoffeeShopRepository: coffeeShopRepo,
	}
}
