package service

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/ifulqt/coffeeshops-api/internal/adapter/cloudflare"
	"github.com/ifulqt/coffeeshops-api/internal/adapter/repository"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/entity"
)

type ImageService interface {
	UploadImages(ctx context.Context, req entity.ImageEntity) error
	UploadImageR2(ctx context.Context, req entity.FileUploadImageEntity) (string, error)
	DeleteImage(ctx context.Context, idCoffeeShop int) error
}

type imageService struct {
	ImageRepository repository.ImageRepository
	R2                    cloudflare.CloudFlareR2Adapter
}

// DeleteImage implements [ImageService].
func (u *imageService) DeleteImage(ctx context.Context, id int) error {
	results, err := u.ImageRepository.GetImageByIDCoffeeShop(ctx, id)
	if err != nil {
		code := "[SERVICE] DeleteCoffeeShop - 1"
		log.Errorw(code, err)
		return err
	}
	for _, res := range results {
		err = u.R2.DeleteImage(ctx, res.Name)
		if err != nil {
			code := "[SERVICE] DeleteCoffeeShop - 2"
			log.Errorw(code, err)
			return err
		}
	}
	return nil
}

// UploadImageR2 implements [UploadImageService].
func (u *imageService) UploadImageR2(ctx context.Context, req entity.FileUploadImageEntity) (string, error) {
	imageURL, err := u.R2.UploadImage(ctx, req)
	if err != nil {
		code := "[SERVICE] UploadImagesR2 - 1"
		log.Errorw(code, err)
		return "", err
	}
	return imageURL, nil
}

// UploadImages implements [UploadImageService].
func (u *imageService) UploadImages(ctx context.Context, req entity.ImageEntity) error {
	err := u.ImageRepository.UploadImages(ctx, req)
	if err != nil {
		code := "[SERVICE] UploadImages - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

func NewImageService(imageRepo repository.ImageRepository, r2 cloudflare.CloudFlareR2Adapter) ImageService {
	return &imageService{
		ImageRepository: imageRepo,
		R2:                    r2,
	}
}
