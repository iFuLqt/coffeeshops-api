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
	UpdateFacilityCoffeeShop(ctx context.Context, req []string, idCoffeeShop int64) error
	CreateFacility(ctx context.Context, req entity.FacilityEntity) error
	UpdateFacility(ctx context.Context, req entity.FacilityEntity, id int64) error
	DeleteFacility(ctx context.Context, id int64) error
	GetFacilities(ctx context.Context) ([]entity.FacilityEntity, error)
}

type facilityRepository struct {
	db *gorm.DB
}

// UpdateFacilityCoffeeShop implements [FacilityRepository].
func (f *facilityRepository) UpdateFacilityCoffeeShop(ctx context.Context, req []string, idCoffeeShop int64) error {
	err := f.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// =====================================================
		// 1. Ambil facility LAMA (facility_id + code)
		// =====================================================
		type oldFacility struct {
			FacilityID int64
			Code       string
		}
		var oldFacilities []oldFacility
		err := tx.Table("coffee_shop_facility csf").Select("csf.facility_id, f.code").Joins("JOIN facilities f ON f.id = csf.facility_id").
			Where("csf.coffee_shop_id = ?", idCoffeeShop).
			Scan(&oldFacilities).Error
		if err != nil {
			code := "[REPOSITORY] UpdateFacilityCoffeeShop - 1"
			log.Errorw(code, err)
			return err
		}
		// =====================================================
		// 2. Map OLD (code -> id)
		// =====================================================
		oldMap := make(map[string]int64)
		for _, f := range oldFacilities {
			oldMap[f.Code] = f.FacilityID
		}
		// =====================================================
		// 3. Map NEW (code)
		// =====================================================
		newMap := make(map[string]struct{})
		for _, code := range req {
			newMap[code] = struct{}{}
		}
		// =====================================================
		// 4. DELETE: old - new
		// =====================================================
		var deleteIDs []int64
		for code, facilityID := range oldMap {
			if _, exists := newMap[code]; !exists {
				deleteIDs = append(deleteIDs, facilityID)
			}
		}
		if len(deleteIDs) > 0 {
			err := tx.Where("coffee_shop_id = ? AND facility_id IN ?", idCoffeeShop, deleteIDs).
				Delete(&model.CoffeeShopFacility{}).Error
			if err != nil {
				code := "[REPOSITORY] UpdateFacilityCoffeeShop - 2"
				log.Errorw(code, err)
				return err
			}
		}
		// =====================================================
		// 5. INSERT: new - old
		// =====================================================
		var insertCodes []string
		for code := range newMap {
			if _, exists := oldMap[code]; !exists {
				insertCodes = append(insertCodes, code)
			}
		}
		if len(insertCodes) == 0 {
			return nil
		}
		var facilities []model.Facility
		err = tx.Where("code IN ?", insertCodes).Find(&facilities).Error
		if err != nil {
			code := "[REPOSITORY] UpdateFacilityCoffeeShop - 3"
			log.Errorw(code, err)
			return err
		}
		var inserts []model.CoffeeShopFacility
		for _, fac := range facilities {
			inserts = append(inserts, model.CoffeeShopFacility{
				CoffeeShopID: idCoffeeShop,
				FacilityID:   fac.ID,
			})
		}
		if len(inserts) > 0 {
			err := tx.Create(&inserts).Error
			if err != nil {
				code := "[REPOSITORY] UpdateFacilityCoffeeShop - 4"
				log.Errorw(code, err)
				return err
			}
		}

		return nil
	})
	if err != nil {
		code := "[REPOSITORY] UpdateFacilityCoffeeShop - 5"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// DeleteFacility implements [FacilityRepository].
func (f *facilityRepository) DeleteFacility(ctx context.Context, id int64) error {
	var modelFacility model.Facility
	result := f.db.WithContext(ctx).Where("id = ?", id).Delete(&modelFacility)
	if result.Error != nil {
		code := "[REPOSITORY] DeleteFacility - 1"
		log.Errorw(code, result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		code := "[REPOSITORY] DeleteFacility - 2"
		log.Errorw(code, domerror.ErrDataNotFound)
		return domerror.ErrDataNotFound
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

func NewFacilityRepository(db *gorm.DB) FacilityRepository {
	return &facilityRepository{
		db: db,
	}
}
