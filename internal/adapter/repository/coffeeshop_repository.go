package repository

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/domerror"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/entity"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/model"
	"gorm.io/gorm"
)

type CoffeeShopRepository interface {
	CreateCoffeeShop(ctx context.Context, req entity.CoffeeShopEntity) (int64, error)
	GetCoffeeShops(ctx context.Context) ([]entity.CoffeeShopEntity, error)
	GetCoffeeShopByID(ctx context.Context, id int64) (*entity.CoffeeShopEntity, error)
	UpdateCoffeeShop(ctx context.Context, req entity.CoffeeShopEntity, idCoffeeShop int64) error
	DeleteCoffeeShop(ctx context.Context, id int64) error
}

type coffeeShopRepository struct {
	db *gorm.DB
}

// CreateCoffeeShop implements [CoffeeShopRepository].
func (c *coffeeShopRepository) CreateCoffeeShop(ctx context.Context, req entity.CoffeeShopEntity) (int64, error) {
	modelCoffe := model.CoffeeShop{
		Name:        req.Name,
		Address:     req.Address,
		Slug: req.Slug,
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

	err := c.db.WithContext(ctx).Create(&modelCoffe).Error
	if err != nil {
		code := "[REPOSITORY] CreateCoffeeShop - 1"
		log.Errorw(code, err)
		return 0, err
	}
	return modelCoffe.ID, nil
}

// DeleteCoffeeShop implements [CoffeeShopRepository].
func (c *coffeeShopRepository) DeleteCoffeeShop(ctx context.Context, id int64) error {
	var modelCoffee model.CoffeeShop
	result := c.db.WithContext(ctx).Where("id = ?", id).Delete(&modelCoffee)
	if result.Error != nil {
		code := "[REPOSITORY] DeleteCoffeeShop - 1"
		log.Errorw(code, result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		code := "[REPOSITORY] DeleteCoffeeShop - 2"
		log.Errorw(code, domerror.ErrDataNotFound)
		return domerror.ErrDataNotFound
	}

	return nil
}

// GetCoffeeShopByID implements [CoffeeShopRepository].
func (c *coffeeShopRepository) GetCoffeeShopByID(ctx context.Context, id int64) (*entity.CoffeeShopEntity, error) {
	var modelCoffe model.CoffeeShop
	err := c.db.WithContext(ctx).Where("id = ?", id).Preload("Category").Preload("UserUpdate").Preload("UserCreate").
		Preload("Images").Preload("CoffeeShopFacility.Facility").First(&modelCoffe).Error
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

	facilities := []entity.FacilityEntity{}
	for _, faci := range modelCoffe.CoffeeShopFacility {
		facility := entity.FacilityEntity{
			Name: faci.Facility.Name,
		}
		facilities = append(facilities, facility)
	}

	coffeeEnt := entity.CoffeeShopEntity{
		ID:        modelCoffe.ID,
		Name:      modelCoffe.Name,
		Slug: modelCoffe.Slug,
		Address:   modelCoffe.Address,
		Latitude:  modelCoffe.Latitude,
		Longitude: modelCoffe.Longitude,
		OpenTime:  modelCoffe.OpenTime,
		CloseTime: modelCoffe.CloseTime,
		Instagram: modelCoffe.Instagram,
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
		Facility: facilities,
		Image:    imageSlics,
	}
	return &coffeeEnt, nil
}

// GetCoffeeShops implements [CoffeeShopRepository].
func (c *coffeeShopRepository) GetCoffeeShops(ctx context.Context) ([]entity.CoffeeShopEntity, error) {
	var modelCoffe []model.CoffeeShop
	err := c.db.WithContext(ctx).Preload("Category").Preload("Images").Find(&modelCoffe).Error
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
			Slug: val.Slug,
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
func (c *coffeeShopRepository) UpdateCoffeeShop(ctx context.Context, req entity.CoffeeShopEntity, idCoffeeShop int64) error {
	modelCoffee := model.CoffeeShop{
		Name: req.Name,
		Slug: req.Slug,
		Address: req.Address,
		Latitude: req.Latitude,
		Longitude: req.Longitude,
		OpenTime: req.OpenTime,
		CloseTime: req.CloseTime,
		Instagram: req.Instagram,
		UpdatedByID: req.UserUpdate.ID,
		CategoryID: req.Category.ID,
		IsActive: req.IsActive,
	}
	result := c.db.Where("id = ?", idCoffeeShop).Updates(modelCoffee)
	if result.Error != nil {
		code := "[REPOSITORY] UpdateCoffeeShop - 1"
		log.Errorw(code, result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		code := "[REPOSITORY] UpdateCoffeeShop - 2"
		log.Errorw(code, domerror.ErrDataNotFound)
		return domerror.ErrDataNotFound
	}
	return nil
}

func NewCoffeeShopRepository(db *gorm.DB) CoffeeShopRepository {
	return &coffeeShopRepository{
		db: db,
	}
}
