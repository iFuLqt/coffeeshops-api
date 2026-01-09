package repository

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/entity"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/model"
	"gorm.io/gorm"
)

type CoffeeShopRepository interface {
	CreateCoffeeShop(ctx context.Context, req entity.CoffeeShopEntity) (int, error)
	GetCoffeeShops(ctx context.Context) ([]entity.CoffeeShopEntity, error)
	GetCoffeeShopByID(ctx context.Context, id int) (*entity.CoffeeShopEntity, error)
	UpdateCoffeeShop(ctx context.Context, req entity.CoffeeShopEntity) error
	DeleteCoffeeShop(ctx context.Context, id int) error
	UploadImages(ctx context.Context, req entity.ImageEntity) error
}

type coffeeShopRepository struct {
	db *gorm.DB
}

// UploadImages implements [CoffeeShopRepository].
func (c *coffeeShopRepository) UploadImages(ctx context.Context, req entity.ImageEntity) error {
	modelCoffe := model.CoffeeShopImage{
		CoffeeShopID: req.CoffeeShopID,
		Image:        req.Image,
		IsPrimary:    req.IsPrimary,
	}
	err := c.db.Create(&modelCoffe).Error
	if err != nil {
		code := "[REPOSITORY] UploadImages - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// CreateCoffeeShop implements [CoffeeShopRepository].
func (c *coffeeShopRepository) CreateCoffeeShop(ctx context.Context, req entity.CoffeeShopEntity) (int, error) {
	modelCoffe := model.CoffeeShop{
		Name:        req.Name,
		Address:     req.Address,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		OpenTime:    req.OpenTimeParse,
		CloseTime:   req.CloseTimeParse,
		Parking:     req.Parking,
		PrayerRoom:  req.PrayerRoom,
		Wifi:        req.Wifi,
		Gofood:      req.Gofood,
		Instagram:   req.Instagram,
		CreatedByID: req.UserCreate.ID,
		UpdatedByID: req.UserUpdate.ID,
		CategoryID:  req.Category.ID,
		IsActive:    req.IsActive,
	}

	err := c.db.Create(&modelCoffe).Error
	if err != nil {
		code := "[REPOSITORY] CreateCoffeeShop - 1"
		log.Errorw(code, err)
		return 0, err
	}
	return modelCoffe.ID, nil
}

// DeleteCoffeeShop implements [CoffeeShopRepository].
func (c *coffeeShopRepository) DeleteCoffeeShop(ctx context.Context, id int) error {
	panic("unimplemented")
}

// GetCoffeeShopByID implements [CoffeeShopRepository].
func (c *coffeeShopRepository) GetCoffeeShopByID(ctx context.Context, id int) (*entity.CoffeeShopEntity, error) {
	panic("unimplemented")
}

// GetCoffeeShops implements [CoffeeShopRepository].
func (c *coffeeShopRepository) GetCoffeeShops(ctx context.Context) ([]entity.CoffeeShopEntity, error) {
	panic("unimplemented")
}

// UpdateCoffeeShop implements [CoffeeShopRepository].
func (c *coffeeShopRepository) UpdateCoffeeShop(ctx context.Context, req entity.CoffeeShopEntity) error {
	panic("unimplemented")
}

func NewCoffeeShopRepository(db *gorm.DB) CoffeeShopRepository {
	return &coffeeShopRepository{
		db: db,
	}
}
