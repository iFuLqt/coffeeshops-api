package repository

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/domerror"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/entity"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/model"
	"gorm.io/gorm"
)

type FacilityRepository interface {
	CreateFacilityCoffeeShop(ctx context.Context, idFacility, idCoffeShop int) error
	CodeForCreateFCS(ctx context.Context, code string) (int, error)
	CreateFacility(ctx context.Context, req entity.FacilityEntity) error
}

type facilityRepository struct {
	db *gorm.DB
}

// CreateFacility implements [FacilityRepository].
func (f *facilityRepository) CreateFacility(ctx context.Context, req entity.FacilityEntity) error {
	var count int64
	err := f.db.Table("facilities").Where("code = ?", req.Code).Count(&count).Error
	if err != nil {
		code := "[REPOSITORY] CreateFacility - 1"
		log.Errorw(code, err)
		return err
	}

	if count > 0 {
		return domerror.ErrDuplicate
	}

	modelFacility := model.Facility{
		Code: req.Code,
		Name: req.Name,
	}
	err = f.db.Create(&modelFacility).Error
	if err != nil {
		code := "[REPOSITORY] CreateFacility - 2"
		log.Errorw(code, err)
		return err
	}
	return nil
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
