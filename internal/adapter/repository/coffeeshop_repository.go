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
		OpenTime:    req.OpenTime,
		CloseTime:   req.CloseTime,
		Instagram:   req.Instagram,
		CreatedByID: req.UserCreate.ID,
		UpdatedByID: req.UserUpdate.ID,
		CategoryID:  req.Category.ID,
		IsActive:    req.IsActive,
		UpdatedAt:   req.UpdatedAt,
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
	var modelCoffe model.CoffeeShop
	err := c.db.Where("id = ?", id).Preload("Category").Preload("UserUpdate").Preload("UserCreate").
		Preload("Images").First(&modelCoffe).Error
	if err != nil {
		code := "[REPOSITORY] GetCoffeeShopByID - 1"
		log.Errorw(code, err)
		return nil, err
	}

	imageSlics := []entity.ImageEntity{}
	for _, valImage := range modelCoffe.Images {
		imageEnt := entity.ImageEntity{
			Image:     valImage.Image,
			IsPrimary: valImage.IsPrimary,
		}
		imageSlics = append(imageSlics, imageEnt)
	}

	coffeeEnt := entity.CoffeeShopEntity{
		ID:         modelCoffe.ID,
		Name:       modelCoffe.Name,
		Address:    modelCoffe.Address,
		Latitude:   modelCoffe.Latitude,
		Longitude:  modelCoffe.Longitude,
		OpenTime:   modelCoffe.OpenTime,
		CloseTime:  modelCoffe.CloseTime,
		Instagram:  modelCoffe.Instagram,
		UserCreate: entity.UserEntity{
			ID:   modelCoffe.UserCreate.ID,
			Name: modelCoffe.UserCreate.Name,
		},
		UserUpdate: entity.UserEntity{
			ID:   modelCoffe.UserUpdate.ID,
			Name: modelCoffe.UserUpdate.Name,
		},
		Category: entity.CategoryEntity{
			ID:   modelCoffe.Category.ID,
			Name: modelCoffe.Category.Name,
		},
		Image: imageSlics,
	}
	return &coffeeEnt, nil
}

// GetCoffeeShops implements [CoffeeShopRepository].
func (c *coffeeShopRepository) GetCoffeeShops(ctx context.Context) ([]entity.CoffeeShopEntity, error) {
	var modelCoffe []model.CoffeeShop
	err := c.db.Preload("Category").Preload("Images").Find(&modelCoffe).Error
	if err != nil {
		code := "[REPOSITORY] GetCoffeeShops - 1"
		log.Errorw(code, err)
		return nil, err
	}
	coffeeEntity := []entity.CoffeeShopEntity{}
	for _, val := range modelCoffe {
		imageSlics := []entity.ImageEntity{}
		for _, valImage := range val.Images {
			if !valImage.IsPrimary {
				continue
			}
			imageEnt := entity.ImageEntity{
				Image:     valImage.Image,
				IsPrimary: true,
			}
			imageSlics = append(imageSlics, imageEnt)
		}

		coffeeEnt := entity.CoffeeShopEntity{
			ID:        val.ID,
			Name:      val.Name,
			Address:   val.Address,
			OpenTime:  val.OpenTime,
			CloseTime: val.CloseTime,
			Instagram: val.Instagram,
			Category: entity.CategoryEntity{
				ID:   val.Category.ID,
				Name: val.Category.Name,
			},
			Image: imageSlics,
		}
		coffeeEntity = append(coffeeEntity, coffeeEnt)
	}
	return coffeeEntity, nil
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
