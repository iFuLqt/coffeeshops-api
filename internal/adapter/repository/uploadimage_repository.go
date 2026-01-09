package repository

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/entity"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/model"
	"gorm.io/gorm"
)

type UploadImageRepository interface {
	UploadImages(ctx context.Context, req entity.ImageEntity) error
}

type uploadImageRepository struct {
	db *gorm.DB
}

// UploadImages implements [UploadImageRepository].
func (u *uploadImageRepository) UploadImages(ctx context.Context, req entity.ImageEntity) error {
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

func NewUploadImageRepository(db *gorm.DB) UploadImageRepository {
	return &uploadImageRepository{
		db: db,
	}
}
