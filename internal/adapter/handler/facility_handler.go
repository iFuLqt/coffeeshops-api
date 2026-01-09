package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/ifulqt/coffeeshops-api/internal/adapter/handler/request"
	"github.com/ifulqt/coffeeshops-api/internal/adapter/handler/response"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/entity"
	"github.com/ifulqt/coffeeshops-api/internal/core/service"
	"github.com/ifulqt/coffeeshops-api/library/helper"
	"github.com/ifulqt/coffeeshops-api/library/validat"
)

type FacilityHandler interface {
	CreateFacilityCoffeeShop(c *fiber.Ctx) error
}

type facilityHandler struct {
	FacilityService service.FacilityService
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

	err = f.FacilityService.CreateFacilityCoffeeShop(c.Context(), req.Code, int(CoffeeShopID))
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
