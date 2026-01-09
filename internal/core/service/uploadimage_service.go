package service

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/ifulqt/coffeeshops-api/internal/adapter/cloudflare"
	"github.com/ifulqt/coffeeshops-api/internal/adapter/repository"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/entity"
)

type UploadImageService interface {
	UploadImages(ctx context.Context, req entity.ImageEntity) error
	UploadImageR2(ctx context.Context, req entity.FileUploadImageEntity) (string, error)
}

type uploadImageService struct {
	UploadImageRepository repository.UploadImageRepository
	R2 cloudflare.CloudFlareR2Adapter
}

// UploadImageR2 implements [UploadImageService].
func (u *uploadImageService) UploadImageR2(ctx context.Context, req entity.FileUploadImageEntity) (string, error) {
	imageURL, err := u.R2.UploadImage(ctx, req)
	if err != nil {
		code := "[SERVICE] UploadImagesR2 - 1"
		log.Errorw(code, err)
		return "", err
	}
	return imageURL, nil
}

// UploadImages implements [UploadImageService].
func (u *uploadImageService) UploadImages(ctx context.Context, req entity.ImageEntity) error {
	err := u.UploadImageRepository.UploadImages(ctx, req)
	if err != nil {
		code := "[SERVICE] UploadImages - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

func NewUploadImageService(uploadImageRepo repository.UploadImageRepository, r2 cloudflare.CloudFlareR2Adapter) UploadImageService {
	return &uploadImageService{
		UploadImageRepository: uploadImageRepo,
		R2: r2,
	}
}
