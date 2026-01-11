package repository

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2/log"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/domerror"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/entity"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/model"
	"gorm.io/gorm"
)

type FacilityRepository interface {
	CreateFacilityCoffeeShop(ctx context.Context, idFacility, idCoffeShop int64) error
	CodeForCreateFCS(ctx context.Context, code string) (int64, error)
	CreateFacility(ctx context.Context, req entity.FacilityEntity) error
	UpdateFacility(ctx context.Context, req entity.FacilityEntity, id int64) error
	DeleteFacility(ctx context.Context, id int64) error
	GetFacilities(ctx context.Context) ([]entity.FacilityEntity, error)
}

type facilityRepository struct {
	db *gorm.DB
}

// DeleteFacility implements [FacilityRepository].
func (f *facilityRepository) DeleteFacility(ctx context.Context, id int64) error {
	var modelFacility model.Facility
	err := f.db.WithContext(ctx).Where("id = ?", id).Delete(&modelFacility).Error
	if err != nil {
		code := "[REPOSITORY] DeleteFacility - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// GetFacilities implements [FacilityRepository].
func (f *facilityRepository) GetFacilities(ctx context.Context) ([]entity.FacilityEntity, error) {
	var modelFac []model.Facility
	err := f.db.WithContext(ctx).Order("id DESC").Find(&modelFac).Error
	if err != nil {
		code := "[REPOSITORY] GetFacilities - 1"
		log.Errorw(code, err)
		return nil, err
	}
	entFacs := []entity.FacilityEntity{}
	for _, val := range modelFac {
		entFac := entity.FacilityEntity{
			ID:   val.ID,
			Code: val.Code,
			Name: val.Name,
		}
		entFacs = append(entFacs, entFac)
	}
	return entFacs, nil
}

// UpdateFacility implements [FacilityRepository].
func (f *facilityRepository) UpdateFacility(ctx context.Context, req entity.FacilityEntity, id int64) error {
	modelFacility := model.Facility{
		Name: req.Name,
		Code: req.Code,
	}
	result := f.db.WithContext(ctx).Where("id = ?", id).Updates(&modelFacility)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			code := "[REPOSITORY] UpdateFacility - 1"
			log.Errorw(code, result.Error)
			return domerror.ErrDuplicate
		}
		code := "[REPOSITORY] UpdateFacility - 2"
		log.Errorw(code, result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		code := "[REPOSITORY] UpdateFacility - 3"
		log.Errorw(code, domerror.ErrDataNotFound)
		return domerror.ErrDataNotFound
	}
	return nil
}

// CreateFacility implements [FacilityRepository].
func (f *facilityRepository) CreateFacility(ctx context.Context, req entity.FacilityEntity) error {
	modelFacility := model.Facility{
		Code: req.Code,
		Name: req.Name,
	}
	err := f.db.WithContext(ctx).Create(&modelFacility).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			code := "[REPOSITORY] CreateFacility - 1"
			log.Errorw(code, err)
			return domerror.ErrDuplicate
		}
		code := "[REPOSITORY] CreateFacility - 2"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// CodeForCreateFCS implements [FacilityRepository].
func (f *facilityRepository) CodeForCreateFCS(ctx context.Context, code string) (int64, error) {
	var modelFacility model.Facility
	err := f.db.WithContext(ctx).Where("code = ?", code).Find(&modelFacility).Error
	if err != nil {
		code := "[REPOSITORY] CodeForCreateFCS - 1"
		log.Errorw(code, err)
		return 0, err
	}
	return modelFacility.ID, nil
}

// CreateFacilityCoffeeShopRelation implements [FacilityRepository].
func (f *facilityRepository) CreateFacilityCoffeeShop(ctx context.Context, idFacility, idCoffeShop int64) error {
	modelCSFacility := model.CoffeeShopFacility{
		CoffeeShopID: idCoffeShop,
		FacilityID:   idFacility,
	}
	err := f.db.WithContext(ctx).Create(&modelCSFacility).Error
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
