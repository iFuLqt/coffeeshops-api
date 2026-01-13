package service

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/ifulqt/coffeeshops-api/internal/adapter/repository"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/entity"
)

type FacilityService interface {
	UpdateFacilityCoffeeShop(ctx context.Context, req []string, idCoffeeShop int64) error
	CreateFacility(ctx context.Context, req entity.FacilityEntity) error
	UpdateFacility(ctx context.Context, req entity.FacilityEntity, id int64) error
	DeleteFacility(ctx context.Context, id int64) error
	GetFacilities(ctx context.Context) ([]entity.FacilityEntity, error)
}

type facilityService struct {
	FacilityRepository repository.FacilityRepository
}

// UpdateFacilityCoffeeShop implements [FacilityService].
func (f *facilityService) UpdateFacilityCoffeeShop(ctx context.Context, req []string, idCoffeeShop int64) error {
	err := f.FacilityRepository.UpdateFacilityCoffeeShop(ctx, req, idCoffeeShop)
	if err != nil {
		code := "[REPOSITORY] UpdateFacilityCoffeeShop - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// DeleteFacility implements [FacilityService].
func (f *facilityService) DeleteFacility(ctx context.Context, id int64) error {
	err := f.FacilityRepository.DeleteFacility(ctx, id)
	if err != nil {
		code := "[SERVICE] DeleteFacility - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// GetFacilities implements [FacilityService].
func (f *facilityService) GetFacilities(ctx context.Context) ([]entity.FacilityEntity, error) {
	results, err := f.FacilityRepository.GetFacilities(ctx)
	if err != nil {
		code := "[SERVICE] GetFacilities - 1"
		log.Errorw(code, err)
		return nil, err
	}
	return results, nil
}

// UpdateFacility implements [FacilityService].
func (f *facilityService) UpdateFacility(ctx context.Context, req entity.FacilityEntity, id int64) error {
	err := f.FacilityRepository.UpdateFacility(ctx, req, id)
	if err != nil {
		code := "[SERVICE] UpdateFacility - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// CreateFacility implements [FacilityService].
func (f *facilityService) CreateFacility(ctx context.Context, req entity.FacilityEntity) error {
	err := f.FacilityRepository.CreateFacility(ctx, req)
	if err != nil {
		code := "[SERVICE] - CreateFacility"
		log.Errorw(code, err)
		return err
	}
	return err
}

func NewFacilityService(facilityRepo repository.FacilityRepository) FacilityService {
	return &facilityService{
		FacilityRepository: facilityRepo,
	}
}
