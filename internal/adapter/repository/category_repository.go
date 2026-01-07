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
		code := "[REPOSITORY] CreateCategory - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// DeleteCategory implements [CategoryRepository].
func (c *categoryRepository) DeleteCategory(ctx context.Context, id int) error {
	var modelCategory model.Category
	err := c.db.Where("id = ?", id).Delete(&modelCategory).Error
	if err != nil {
		code := "[REPOSITORY] DeleteCategory - 2"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// GetCategories implements [CategoryRepository].
func (c *categoryRepository) GetCategories(ctx context.Context) ([]entity.CategoryEntity, error) {
	var modelCategory []model.Category
	err := c.db.Order("created_at DESC").Preload("User").Find(&modelCategory).Error
	if err != nil {
		code := "[REPOSITORY] GetCategories - 1"
		log.Errorw(code, err)
		return nil, err
	}

	resps := []entity.CategoryEntity{}
	for _, val := range modelCategory {
		resps = append(resps, entity.CategoryEntity{
			ID:   val.ID,
			Name: val.Name,
			Slug: val.Slug,
			User: entity.UserEntity{
				ID:   val.User.ID,
				Name: val.User.Name,
			},
		})
	}

	return resps, nil

}

// GetCategoryByID implements [CategoryRepository].
func (c *categoryRepository) GetCategoryByID(ctx context.Context, id int) (*entity.CategoryEntity, error) {
	var categoryModel model.Category
	err := c.db.Where("id = ?", id).Preload("User").First(&categoryModel).Error
	if err != nil {
		code := "[REPOSITORY] GetCategoryByID - 1"
		log.Errorw(code, err)
		return nil, err
	}

	resp := entity.CategoryEntity{
		ID:   categoryModel.ID,
		Name: categoryModel.Name,
		Slug: categoryModel.Slug,
		User: entity.UserEntity{
			ID:    categoryModel.User.ID,
			Name:  categoryModel.User.Name,
			Email: categoryModel.User.Email,
		},
	}

	return &resp, nil
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
