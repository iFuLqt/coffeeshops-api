package service

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/ifulqt/coffeeshops-api/internal/adapter/repository"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/entity"
	"github.com/ifulqt/coffeeshops-api/library/helper"
)

type CategoryService interface {
	GetCategories(ctx context.Context) ([]entity.CategoryEntity, error)
	GetCategoryByID(ctx context.Context, id int64) (*entity.CategoryEntity, error)
	CreateCategory(ctx context.Context, req entity.CategoryEntity) error
	UpdateCategory(ctx context.Context, req entity.CategoryEntity) error
	DeleteCategory(ctx context.Context, id int64) error
}

type categoryService struct {
	CategoryRepository repository.CategoryRepository
}

// CreateCategory implements [CategoryService].
func (c *categoryService) CreateCategory(ctx context.Context, req entity.CategoryEntity) error {
	slug := helper.GenerateSlug(req.Name)
	req.Slug = slug

	err := c.CategoryRepository.CreateCategory(ctx, req)
	if err != nil {
		code := "[SERVICE] CreateCategory - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// DeleteCategory implements [CategoryService].
func (c *categoryService) DeleteCategory(ctx context.Context, id int64) error {
	err := c.CategoryRepository.DeleteCategory(ctx, id)
	if err != nil {
		code := "[SERVICE] DeleteCategory - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// GetCategories implements [CategoryService].
func (c *categoryService) GetCategories(ctx context.Context) ([]entity.CategoryEntity, error) {
	results, err := c.CategoryRepository.GetCategories(ctx)
	if err != nil {
		code := "[SERVICE] GetCategories - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return results, err
}

// GetCategoryByID implements [CategoryService].
func (c *categoryService) GetCategoryByID(ctx context.Context, id int64) (*entity.CategoryEntity, error) {
	result, err := c.CategoryRepository.GetCategoryByID(ctx, id)
	if err != nil {
		code := "[SERVICE] GetCategoryByID - 1"
		log.Errorw(code, err)
		return nil, err
	}
	return result, nil
}

// UpdateCategory implements [CategoryService].
func (c *categoryService) UpdateCategory(ctx context.Context, req entity.CategoryEntity) error {
	slug := helper.GenerateSlug(req.Name)
	req.Slug = slug

	err := c.CategoryRepository.UpdateCategory(ctx, req)
	if err != nil {
		code := "[SERVICE] UpdateCategory - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{
		CategoryRepository: categoryRepo,
	}
}
