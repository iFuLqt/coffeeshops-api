package service

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/ifulqt/coffeeshops-api/internal/adapter/cloudflare"
	"github.com/ifulqt/coffeeshops-api/internal/adapter/repository"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/entity"
)

type CoffeeShopService interface {
	CreateCoffeeShop(ctx context.Context, req entity.CoffeeShopEntity) (int, error)
	GetCoffeeShops(ctx context.Context) ([]entity.CoffeeShopEntity, error)
	GetCoffeeShopByID(ctx context.Context, id int) (*entity.CoffeeShopEntity, error)
	UpdateCoffeeShop(ctx context.Context, req entity.CoffeeShopEntity) error
	DeleteCoffeeShop(ctx context.Context, id int) error
	UploadImages(ctx context.Context, req entity.ImageEntity) error
	UploadImagesR2(ctx context.Context, req entity.FileUploadImageEntity) (string, error)
}

type coffeeShopService struct {
	CoffeeShopRepository repository.CoffeeShopRepository
	R2                   cloudflare.CloudFlareR2Adapter
}

// UploadImagesR2 implements [CoffeeShopService].
func (c *coffeeShopService) UploadImagesR2(ctx context.Context, req entity.FileUploadImageEntity) (string, error) {
	imageURL, err := c.R2.UploadImage(ctx, req)
	if err != nil {
		code := "[SERVICE] UploadImagesR2 - 1"
		log.Errorw(code, err)
		return "", err
	}
	return imageURL, nil
}

// UploadImages implements [CoffeeShopService].
func (c *coffeeShopService) UploadImages(ctx context.Context, req entity.ImageEntity) error {
	err := c.CoffeeShopRepository.UploadImages(ctx, req)
	if err != nil {
		code := "[SERVICE] UploadImages - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// CreateCoffeeShop implements [CoffeeShopService].
func (c *coffeeShopService) CreateCoffeeShop(ctx context.Context, req entity.CoffeeShopEntity) (int, error) {
	id, err := c.CoffeeShopRepository.CreateCoffeeShop(ctx, req)
	if err != nil {
		code := "[SERVICE] CreateCoffeeShop - 3"
		log.Errorw(code, err)
		return 0, err
	}
	return id, nil
}

// DeleteCoffeeShop implements [CoffeeShopService].
func (c *coffeeShopService) DeleteCoffeeShop(ctx context.Context, id int) error {
	panic("unimplemented")
}

// GetCoffeeShopByID implements [CoffeeShopService].
func (c *coffeeShopService) GetCoffeeShopByID(ctx context.Context, id int) (*entity.CoffeeShopEntity, error) {
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
func (c *coffeeShopService) UpdateCoffeeShop(ctx context.Context, req entity.CoffeeShopEntity) error {
	panic("unimplemented")
}

func NewCoffeeShopService(coffeeShopRepo repository.CoffeeShopRepository, R2 cloudflare.CloudFlareR2Adapter) CoffeeShopService {
	return &coffeeShopService{
		CoffeeShopRepository: coffeeShopRepo,
		R2:                   R2,
	}
}
