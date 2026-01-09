package cloudflare

import (
	"context"
	"fmt"
	"mime"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2/log"

	"github.com/ifulqt/coffeeshops-api/config"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/entity"
)

type CloudFlareR2Adapter interface {
	UploadImage(ctx context.Context, req entity.FileUploadImageEntity) (string, error)
}

type cloudFlareR2Adapter struct {
	client  *s3.Client
	bucket  string
	baseURL string
}

func (c *cloudFlareR2Adapter) UploadImage(ctx context.Context, req entity.FileUploadImageEntity) (string, error) {
	file, err := req.File.Open()
	if err != nil {
		code := "[CLOUDFLARE] UploadImage - 1"
		log.Errorw(code, err)
		return "", err
	}
	defer file.Close()

	ext := filepath.Ext(req.Name)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "image/jpeg"
	}

	_, err = c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(req.Name),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		code := "[CLOUDFLARE] UploadImage - 2"
		log.Errorw(code, err)
		return "", err
	}

	return fmt.Sprintf("%s/%s", c.baseURL, req.Name), nil
}

func NewCloudFlareR2Adapter(client *s3.Client, cfg *config.Config) CloudFlareR2Adapter {
	s3Client := s3.NewFromConfig(cfg.LoadAwsConfig(), func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.R2.AccountID))
	})

	return &cloudFlareR2Adapter{
		client:  s3Client,
		bucket:  cfg.R2.Name,
		baseURL: cfg.R2.PublicURL,
	}
}
