package repository

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/domerror"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/entity"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/model"
	"gorm.io/gorm"
)

type ImageRepository interface {
	UploadImages(ctx context.Context, req entity.ImageEntity) error
	GetImageByIDCoffeeShop(ctx context.Context, idCoffeeShop int64) ([]entity.FileDeleteImageEntity, error)
	GetImageByIDs(ctx context.Context, idImage []int64, idCoffeeShop int64) ([]entity.ImageEntity, error)
	DeleteImagesForCoffeeShop(ctx context.Context, idImage []int64, idCoffeeShop int64) error
	UpdatePrimaryImage(ctx context.Context, idImage int64, idCoffeeShop int64) error
}

type imageRepository struct {
	db *gorm.DB
}

// UpdatePrimaryImage implements [ImageRepository].
func (u *imageRepository) UpdatePrimaryImage(ctx context.Context, idImage int64, idCoffeeShop int64) error {
	var modelImage model.CoffeeShopImage
	err := u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&modelImage).Where("id = ?", idImage).Where("coffee_shop_id = ?", idCoffeeShop).
			Update("is_primary", true)
		if result.Error != nil {
			code := "[REPOSITORY] UpdatePrimaryImage - 1"
			log.Errorw(code, result.Error)
			return result.Error
		}

		if result.RowsAffected == 0 {
			code := "[REPOSITORY] UpdatePrimaryImage - 2"
			log.Errorw(code, domerror.ErrDataNotFound)
			return domerror.ErrDataNotFound
		}

		err := tx.Model(&modelImage).Where("coffee_shop_id = ? AND id != ?", idCoffeeShop, idImage).
			Update("is_primary", false).Error
		if err != nil {
			code := "[REPOSITORY] UpdatePrimaryImage - 3"
			log.Errorw(code, err)
			return err
		}
		return nil
	})
	if err != nil {
		code := "[REPOSITORY] UpdatePrimaryImage - 4"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// DeleteImageCoffeeShopByIDImage implements [ImageRepository].
func (u *imageRepository) DeleteImagesForCoffeeShop(ctx context.Context, idImage []int64, idCoffeeShop int64) error {
	var modelImage model.CoffeeShopImage
	err := u.db.Where("id IN ?", idImage).Where("coffee_shop_id = ?", idCoffeeShop).Delete(&modelImage).Error
	if err != nil {
		code := "[REPOSITORY] DeleteImagesForCoffeeShop - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// GetImageByIDp implements [ImageRepository].
func (u *imageRepository) GetImageByIDs(ctx context.Context, idImage []int64, idCoffeeShop int64) ([]entity.ImageEntity, error) {
	var modelImage []model.CoffeeShopImage
	err := u.db.Where("id IN ?", idImage).Where("coffee_shop_id = ?", idCoffeeShop).Find(&modelImage).Error
	if err != nil {
		code := "[REPOSITORY] GetImageByID - 1"
		log.Errorw(code, err)
		return nil, err
	}
	resEntitys := []entity.ImageEntity{}
	for _, res := range modelImage {
		resEntity := entity.ImageEntity{
			Image: res.Image,
		}
		resEntitys = append(resEntitys, resEntity)
	}
	return resEntitys, nil
}

// GetImageByIDCoffeeShop implements [UploadImageRepository].
func (u *imageRepository) GetImageByIDCoffeeShop(ctx context.Context, idCoffeeShop int64) ([]entity.FileDeleteImageEntity, error) {
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
