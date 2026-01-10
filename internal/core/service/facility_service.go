package service

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/ifulqt/coffeeshops-api/internal/adapter/repository"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/entity"
)

type FacilityService interface {
	CreateFacilityCoffeeShop(ctx context.Context, req []string, id int) error
	CreateFacility(ctx context.Context, req entity.FacilityEntity) error
}

type facilityService struct {
	FacilityRepository repository.FacilityRepository
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

// CreateFacilityCoffeeShop implements [FacilityService].
func (f *facilityService) CreateFacilityCoffeeShop(ctx context.Context, facilityCodes []string, idCoffeeShop int) error {
	for _, code := range facilityCodes {
		idFacility, err := f.FacilityRepository.CodeForCreateFCS(ctx, code)
		if err != nil {
			code := "[SERVICE] CreateFacilityCoffeeShop - 1"
			log.Errorw(code, err)
			return err
		}
		err = f.FacilityRepository.CreateFacilityCoffeeShop(ctx, idFacility, idCoffeeShop)
		if err != nil {
			code := "[SERVICE] CreateFacilityCoffeeShop - 2"
			log.Errorw(code, err)
			return err
		}
	}
	return nil
}

func NewFacilityService(facilityRepo repository.FacilityRepository) FacilityService {
	return &facilityService{
		FacilityRepository: facilityRepo,
	}
}
