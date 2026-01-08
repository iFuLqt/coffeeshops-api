package repository

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/domerror"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/entity"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/model"
	"github.com/ifulqt/coffeeshops-api/library/helper"
	"gorm.io/gorm"
)

type CoffeeShopRepository interface {
	CreateCoffeeShop(ctx context.Context, req entity.CoffeeShopEntity) error
	GetCoffeeShops(ctx context.Context) ([]entity.CoffeeShopEntity, error)
	GetCoffeeShopByID(ctx context.Context, id int) (*entity.CoffeeShopEntity, error)
	UpdateCoffeeShop(ctx context.Context, req entity.CoffeeShopEntity) error
	DeleteCoffeeShop(ctx context.Context, id int) error
}

type coffeeShopRepository struct {
	db *gorm.DB
}

// CreateCoffeeShop implements [CoffeeShopRepository].
func (c *coffeeShopRepository) CreateCoffeeShop(ctx context.Context, req entity.CoffeeShopEntity) error {
	parsOpen, err := helper.ParseHourMinute(req.OpenTime)
	if err != nil {
		code := "[HANDLER] CreateCoffeeShop - 1"
		log.Errorw(code, err)
		return domerror.ErrParsingTime
	}
	parsClose, err := helper.ParseHourMinute(req.CloseTime)
	if err != nil {
		code := "[HANDLER] CreateCoffeeShop - 1"
		log.Errorw(code, err)
		return domerror.ErrParsingTime
	}
	modelCoffe := model.CoffeeShop{
		Name:        req.Name,
		Address:     req.Address,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		OpenTime:    parsOpen,
		CloseTime:   parsClose,
		Parking:     req.Parking,
		PrayerRoom:  req.PrayerRoom,
		Wifi:        req.Wifi,
		Gofood:      req.Gofood,
		Instagram:   req.Instagram,
		CreatedByID: req.UserCreate.ID,
		UpdatedByID: req.UserUpdate.ID,
		CategoryID:  req.Category.ID,
		IsActive: req.IsActive,
	}

	err = c.db.Create(&modelCoffe).Error
	if err != nil {
		code := "[REPOSITORY] CreateCoffeeShop - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
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
