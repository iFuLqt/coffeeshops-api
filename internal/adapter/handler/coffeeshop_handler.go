package handler

import (
	"errors"
	"fmt"
	"time"

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

type CoffeeShopHandler interface {
	CreateCoffeeShop(c *fiber.Ctx) error
	GetCoffeeShops(c *fiber.Ctx) error
	GetCoffeeShopByID(c *fiber.Ctx) error
	UpdateCoffeeShop(c *fiber.Ctx) error
	DeleteCoffeeShop(c *fiber.Ctx) error
	UploadImages(c *fiber.Ctx) error
}

type coffeeShopHandler struct {
	CoffeeShopService service.CoffeeShopService
}

// UploadImage implements [CoffeeShopHandler].
func (f *coffeeShopHandler) UploadImages(c *fiber.Ctx) error {
	var errResp response.DefaultErrorResponse
	var resp response.DefaultSuccessResponse

	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		errResp.Meta.Status = false
		errResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errResp)
	}

	param := c.Params("coffeeshopID")
	coffeeShopID, err := helper.StringToInt(param)
	if err != nil || coffeeShopID <= 0 {
		code := "[HANDLER] UploadImage - 1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Invalid coffee shop id"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	form, err := c.MultipartForm()
	if err != nil {
		code := "[HANDLER] UploadImage - 2"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Invalid multipart form"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	files := form.File["images"]
	if len(files) == 0 {
		errResp.Meta.Status = false
		errResp.Meta.Message = "Images are required"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	var imageURLs []string

	for i, file := range files {
		uploadReq := entity.FileUploadImageEntity{
			File: file,
			Name: fmt.Sprintf("%d-%d", int(userID), time.Now().Unix()),
		}

		imageURL, err := f.CoffeeShopService.UploadImagesR2(c.Context(), uploadReq)
		if err != nil {
			code := "[HANDLER]"
			log.Errorw(code, err)
			errResp.Meta.Status = false
			errResp.Meta.Message = "Failed to upload image"
			errResp.Meta.Errors = nil
			return c.Status(fiber.StatusInternalServerError).JSON(errResp)
		}

		reqEntity := entity.ImageEntity{
			CoffeeShopID: int(coffeeShopID),
			Image:        imageURL,
			IsPrimary:    i == 0,
		}

		err = f.CoffeeShopService.UploadImages(c.Context(), reqEntity)
		if err != nil {
			code := "[HANDLER] UploadImage - 1"
			log.Errorw(code, err)
			errResp.Meta.Status = false
			errResp.Meta.Message = "Failed to save image"
			errResp.Meta.Errors = nil
			return c.Status(fiber.StatusInternalServerError).JSON(errResp)
		}

		imageURLs = append(imageURLs, imageURL)
	}

	resp.Meta.Status = true
	resp.Meta.Message = "Upload images successfully"
	resp.Meta.Errors = nil
	resp.Data = fiber.Map{
		"images": imageURLs,
	}

	return c.JSON(resp)
}

// CreateCoffeeShop implements [CoffeeShopHandler].
func (f *coffeeShopHandler) CreateCoffeeShop(c *fiber.Ctx) error {
	var errResp response.DefaultErrorResponse
	var resp response.DefaultSuccessResponse
	var req request.CoffeeShopRequest

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
		code := "[HANDLER] CreateCoffeeShop - 1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Invalid request body"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	err = validat.ValidateStruct(&req)
	if err != nil {
		code := "[HANDLER] CreateCoffeeShop - 2"
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

	reqEntity := entity.CoffeeShopEntity{
		Name:       req.CoffeeShop,
		Address:    req.Address,
		Latitude:   req.Latitude,
		Longitude:  req.Longitude,
		OpenTime:   req.OpenTime,
		CloseTime:  req.CloseTime,
		Parking:    req.Parking,
		PrayerRoom: req.PrayerRoom,
		Wifi:       req.Wifi,
		Gofood:     req.Gofood,
		Instagram:  req.Instagram,
		UserCreate: entity.UserEntity{
			ID: int(userID),
		},
		UserUpdate: entity.UserEntity{
			ID: int(userID),
		},
		Category: entity.CategoryEntity{
			ID: req.CategoryID,
		},
		IsActive: true,
	}

	id, err := f.CoffeeShopService.CreateCoffeeShop(c.Context(), reqEntity)
	if err != nil {
		code := "[HANDLER] CreateCoffeeShop - 3"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Errors = nil
		if errors.Is(err, domerror.ErrParsingTime) {
			errResp.Meta.Message = "Invalid request body"
			return c.Status(fiber.StatusBadRequest).JSON(errResp)
		}
		errResp.Meta.Message = "Internal server error"
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	respID := response.CreateCoffeeShopResponse{
		ID: id,
	}

	resp.Meta.Status = true
	resp.Meta.Message = "Create coffee shop successfully"
	resp.Meta.Errors = nil
	resp.Data = respID

	return c.JSON(resp)
}

// DeleteCoffeeShop implements [CoffeeShopHandler].
func (f *coffeeShopHandler) DeleteCoffeeShop(c *fiber.Ctx) error {
	panic("unimplemented")
}

// GetCoffeeShopByID implements [CoffeeShopHandler].
func (f *coffeeShopHandler) GetCoffeeShopByID(c *fiber.Ctx) error {
	panic("unimplemented")
}

// GetCoffeeShops implements [CoffeeShopHandler].
func (f *coffeeShopHandler) GetCoffeeShops(c *fiber.Ctx) error {
	panic("unimplemented")
}

// UpdateCoffeeShop implements [CoffeeShopHandler].
func (f *coffeeShopHandler) UpdateCoffeeShop(c *fiber.Ctx) error {
	panic("unimplemented")
}

func NewCoffeeShopHandler(coffeeShopServ service.CoffeeShopService) CoffeeShopHandler {
	return &coffeeShopHandler{
		CoffeeShopService: coffeeShopServ,
	}
}
