package repository

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/entity"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/model"
	"gorm.io/gorm"
)

type ImageRepository interface {
	UploadImages(ctx context.Context, req entity.ImageEntity) error
	GetImageByIDCoffeeShop(ctx context.Context, idCoffeeShop int) ([]entity.FileDeleteImageEntity, error)
}

type imageRepository struct {
	db *gorm.DB
}

// GetImageByIDCoffeeShop implements [UploadImageRepository].
func (u *imageRepository) GetImageByIDCoffeeShop(ctx context.Context, idCoffeeShop int) ([]entity.FileDeleteImageEntity, error) {
	var modelImage []model.CoffeeShopImage
	err := u.db.Where("coffee_shop_id = ?", idCoffeeShop).Find(&modelImage).Error
	if err != nil {
		code := "[HANDLER] GetImageByIDCoffeeShop - 1"
		log.Errorw(code, err)
		return nil, err
	}

	resEntitys := []entity.FileDeleteImageEntity{}
	for _, val := range modelImage {
		resEntity := entity.FileDeleteImageEntity{
			Name: val.Image,
		}
		resEntitys = append(resEntitys, resEntity)
	}

	return resEntitys, nil
}

// UploadImages implements [UploadImageRepository].
func (u *imageRepository) UploadImages(ctx context.Context, req entity.ImageEntity) error {
	modelCoffe := model.CoffeeShopImage{
		CoffeeShopID: req.CoffeeShopID,
		Image:        req.Image,
		IsPrimary:    req.IsPrimary,
	}
	err := u.db.Create(&modelCoffe).Error
	if err != nil {
		code := "[REPOSITORY] UploadImages - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

func NewImageRepository(db *gorm.DB) ImageRepository {
	return &imageRepository{
		db: db,
	}
}
