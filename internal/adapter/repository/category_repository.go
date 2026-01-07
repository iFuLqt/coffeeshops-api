package repository

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/entity"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/model"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategories(ctx context.Context) ([]entity.CategoryEntity, error)
	GetCategoryByID(ctx context.Context, id int) (*entity.CategoryEntity, error)
	CreateCategory(ctx context.Context, req entity.CategoryEntity) error
	UpdateCategory(ctx context.Context, req entity.CategoryEntity) error
	DeleteCategory(ctx context.Context, id int) error
}

type categoryRepository struct {
	db *gorm.DB
}

// CreateCategory implements [CategoryRepository].
func (c *categoryRepository) CreateCategory(ctx context.Context, req entity.CategoryEntity) error {
	modelCategory := model.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		CreatedByID: req.User.ID,
	}
	err := c.db.Create(&modelCategory).Error
	if err != nil {
		code := "[SERVICE] CreateCategory - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// DeleteCategory implements [CategoryRepository].
func (c *categoryRepository) DeleteCategory(ctx context.Context, id int) error {
	panic("unimplemented")
}

// GetCategories implements [CategoryRepository].
func (c *categoryRepository) GetCategories(ctx context.Context) ([]entity.CategoryEntity, error) {
	panic("unimplemented")
}

// GetCategoryByID implements [CategoryRepository].
func (c *categoryRepository) GetCategoryByID(ctx context.Context, id int) (*entity.CategoryEntity, error) {
	panic("unimplemented")
}

// UpdateCategory implements [CategoryRepository].
func (c *categoryRepository) UpdateCategory(ctx context.Context, req entity.CategoryEntity) error {
	panic("unimplemented")
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}
