package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/ifulqt/coffeeshops-api/internal/adapter/handler/request"
	"github.com/ifulqt/coffeeshops-api/internal/adapter/handler/response"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/domerror"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/entity"
	"github.com/ifulqt/coffeeshops-api/internal/core/service"
	"github.com/ifulqt/coffeeshops-api/library/helper"
	"github.com/ifulqt/coffeeshops-api/library/validat"
)

type FacilityHandler interface {
	CreateFacilityCoffeeShop(c *fiber.Ctx) error
	CreateFacility(c *fiber.Ctx) error
	UpdateFacility(c *fiber.Ctx) error
	DeleteFacility(c *fiber.Ctx) error
	GetFacilities(c *fiber.Ctx) error
}

type facilityHandler struct {
	FacilityService service.FacilityService
}

// DeleteFacility implements [FacilityHandler].
func (f *facilityHandler) DeleteFacility(c *fiber.Ctx) error {
	var resp response.DefaultSuccessResponse
	var errResp response.DefaultErrorResponse

	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		errResp.Meta.Status = false
		errResp.Meta.Message = "Unauthorizes access"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusUnauthorized).JSON(errResp)
	}

	idParameter := c.Params("facilityID")
	idFacility, err := helper.StringToInt(idParameter)
	if err != nil {
		code := "[HANDLER] DeleteFacility - 1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Facility ID must be integer"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	err = f.FacilityService.DeleteFacility(c.Context(), int64(idFacility))
	if err != nil {
		code := "[HANDLER] DeleteFacility - 2"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Errors = nil
		if errors.Is(err, domerror.ErrDataNotFound) {
			errResp.Meta.Message = "Facility not found"
			return c.Status(fiber.StatusNotFound).JSON(errResp)
		}
		errResp.Meta.Message = "Internal server error"
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	resp.Meta.Status = true
	resp.Meta.Message = "Delete facility successfully"
	resp.Meta.Errors = nil

	return c.JSON(resp)
}

// GetFacilities implements [FacilityHandler].
func (f *facilityHandler) GetFacilities(c *fiber.Ctx) error {
	var resp response.DefaultSuccessResponse
	var errResp response.DefaultErrorResponse

	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		errResp.Meta.Status = false
		errResp.Meta.Message = "Unauthorized access"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusUnauthorized).JSON(errResp)
	}

	results, err := f.FacilityService.GetFacilities(c.Context())
	if err != nil {
		code := "[HANDLER] GetFacilities - 1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Internal server error"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	respDatas := []response.FacilityResponse{}
	for _, res := range results {
		respData := response.FacilityResponse{
			ID:   res.ID,
			Name: res.Name,
			Code: res.Code,
		}
		respDatas = append(respDatas, respData)
	}

	resp.Meta.Status = true
	resp.Meta.Message = "Successfully fetched facilities"
	resp.Meta.Errors = nil
	resp.Data = respDatas

	return c.JSON(resp)
}

// UpdateFacility implements [FacilityHandler].
func (f *facilityHandler) UpdateFacility(c *fiber.Ctx) error {
	var resp response.DefaultSuccessResponse
	var errResp response.DefaultErrorResponse
	var req request.FacilityRequest

	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		errResp.Meta.Status = false
		errResp.Meta.Message = "Unauthorized access"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusUnauthorized).JSON(errResp)
	}

	idParameter := c.Params("facilityID")
	idFacility, err := helper.StringToInt(idParameter)
	if err != nil {
		code := "[HANDLER] UpdateFacility - 1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Facility ID must be integer"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	err = c.BodyParser(&req)
	if err != nil {
		code := "[HANDLER] UpdateFacility - 2"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Invalid request body"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	err = validat.ValidateStruct(&req)
	if err != nil {
		code := "[HANDLER] UpdateFacility - 3"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Invalid request data"
		val, ok := err.(validat.ValidationError)
		if ok {
			errResp.Meta.Errors = val.Message
		} else {
			errResp.Meta.Errors = nil
		}
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	reqEntity := entity.FacilityEntity{
		Name: req.Name,
		Code: req.Code,
	}

	err = f.FacilityService.UpdateFacility(c.Context(), reqEntity, int64(idFacility))
	if err != nil {
		code := "[HANDLER] UpdateFacility - 4"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Errors = nil
		if errors.Is(err, domerror.ErrDataNotFound) {
			errResp.Meta.Message = "Facility not found"
			return c.Status(fiber.StatusNotFound).JSON(errResp)
		} else if errors.Is(err, domerror.ErrDuplicate) {
			errResp.Meta.Message = "Code is not ready"
			return c.Status(fiber.StatusBadRequest).JSON(errResp)
		}
		errResp.Meta.Message = "Internal server error"
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	resp.Meta.Status = true
	resp.Meta.Message = "Update facility successfully"
	resp.Meta.Errors = nil

	return c.JSON(resp)
}

// CreateFacility implements [FacilityHandler].
func (f *facilityHandler) CreateFacility(c *fiber.Ctx) error {
	var resp response.DefaultSuccessResponse
	var errResp response.DefaultErrorResponse
	var req request.FacilityRequest

	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		errResp.Meta.Status = false
		errResp.Meta.Message = "Unauthorized access"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusUnauthorized).JSON(errResp)
	}

	err := c.BodyParser(&req)
	if err != nil {
		code := "[HANDLER] CreateFacility - 1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Invalid request body"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	err = validat.ValidateStruct(&req)
	if err != nil {
		code := "[HANDLER] CreateFacility - 2"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Invalid request data"
		val, ok := err.(validat.ValidationError)
		if ok {
			errResp.Meta.Errors = val.Message
		} else {
			errResp.Meta.Errors = nil
		}
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	reqEntity := entity.FacilityEntity{
		Code: req.Code,
		Name: req.Name,
	}

	err = f.FacilityService.CreateFacility(c.Context(), reqEntity)
	if err != nil {
		code := "[HANDLER] CreateFacility - 3"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Errors = nil
		if errors.Is(err, domerror.ErrDuplicate) {
			errResp.Meta.Message = "Code is not ready"
			return c.Status(fiber.StatusBadRequest).JSON(errResp)
		} else {
			errResp.Meta.Message = "Internal server error"
			return c.Status(fiber.StatusInternalServerError).JSON(errResp)
		}
	}

	resp.Meta.Status = true
	resp.Meta.Message = "Create facility successfully"
	resp.Meta.Errors = nil

	return c.JSON(resp)
}

// CreateFacilityCoffeeShop implements [FacilityHandler].
func (f *facilityHandler) CreateFacilityCoffeeShop(c *fiber.Ctx) error {
	var resp response.DefaultSuccessResponse
	var errResp response.DefaultErrorResponse
	var req request.FacilityCoffeeShopRequest

	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		errResp.Meta.Status = false
		errResp.Meta.Message = "Unauthorized access"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusUnauthorized).JSON(errResp)
	}

	idParameter := c.Params("coffeeshopID")
	CoffeeShopID, err := helper.StringToInt(idParameter)
	if err != nil {
		code := "[HANDLER] CreateFacilityCoffeeShop - 1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Coffee shop id salah"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	err = c.BodyParser(&req)
	if err != nil {
		code := "[HANDLER] CreateFacilityCoffeeShop - 2"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Invalid request body"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	// if len(req.Code) == 0 {
	// 	errResp.Meta.Status = false
	// 	errResp.Meta.Message = "Facility code cannot be empty"
	// 	errResp.Meta.Errors = nil
	// 	return c.Status(fiber.StatusBadRequest).JSON(errResp)
	// }

	err = validat.ValidateStruct(&req)
	if err != nil {
		code := "[HANDLER] CreateFacilityCoffeeShop - e"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Invalid request data"
		val, ok := err.(validat.ValidationError)
		if ok {
			errResp.Meta.Errors = val.Message
		} else {
			errResp.Meta.Errors = nil
		}
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	err = f.FacilityService.CreateFacilityCoffeeShop(c.Context(), req.FacilityCode, int64(CoffeeShopID))
	if err != nil {
		code := "[HANDLER] CreateFacilituCoffeeShop - 1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Internal server error"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	resp.Meta.Status = true
	resp.Meta.Message = "Successfully create relation facility coffeeshop"
	resp.Meta.Errors = nil

	return c.JSON(resp)
}

func NewFacilityHandler(facilityServ service.FacilityService) FacilityHandler {
	return &facilityHandler{
		FacilityService: facilityServ,
	}
}
