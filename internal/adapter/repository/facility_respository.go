package repository

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/model"
	"gorm.io/gorm"
)

type FacilityRepository interface {
	CreateFacilityCoffeeShop(ctx context.Context, idFacility, idCoffeShop int) error
	CodeForCreateFCS(ctx context.Context, code string) (int, error)
}

type facilityRepository struct {
	db *gorm.DB
}

// CodeForCreateFCS implements [FacilityRepository].
func (f *facilityRepository) CodeForCreateFCS(ctx context.Context, code string) (int, error) {
	var modelFacility model.Facility
	err := f.db.Where("code = ?", code).Find(&modelFacility).Error
	if err != nil {
		code := "[REPOSITORY] CodeForCreateFCS - 1"
		log.Errorw(code, err)
		return 0, err
	}
	return modelFacility.ID, nil
}

// CreateFacilityCoffeeShopRelation implements [FacilityRepository].
func (f *facilityRepository) CreateFacilityCoffeeShop(ctx context.Context, idFacility, idCoffeShop int) error {
	modelCSFacility := model.CoffeeShopFacility{
		CoffeeShopID: idCoffeShop,
		FacilityID:   idFacility,
	}
	err := f.db.Create(&modelCSFacility).Error
	if err != nil {
		code := "[REPOSITORY] CreateFacilityCoffeeShop - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

func NewFacilityRepository(db *gorm.DB) FacilityRepository {
	return &facilityRepository{
		db: db,
	}
}
