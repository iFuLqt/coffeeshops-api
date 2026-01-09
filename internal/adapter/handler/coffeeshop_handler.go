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

	idParam := c.Params("coffeeshopID")
	id, err := helper.StringToInt(idParam)
	if err != nil {
		errResp.Meta.Status = false
		errResp.Meta.Message = "Coffee shop ID salah"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	result, err := f.CoffeeShopService.GetCoffeeShopByID(c.Context(), int(id))
	if err != nil {
		errResp.Meta.Status = false
		errResp.Meta.Message = "Coffee shop ID salah"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	images := []response.ImagesCoffeeShopResponse{}
	for _, res := range result.Image {
		imag := response.ImagesCoffeeShopResponse{
			Image:     res.Image,
			IsPrimary: res.IsPrimary,
		}
		images = append(images, imag)
	}

	respData := response.CoffeeShopByIDResponse{
		ID:        result.ID,
		Name:      result.Name,
		Address:   result.Address,
		OpenClose: helper.GenerateOpenTime(result.OpenTime, result.CloseTime),
		Facility: &response.FacilityCoffeeShopResponse{
			Parking:    result.Parking,
			PrayerRoom: result.PrayerRoom,
			Wifi:       result.Wifi,
			Gofood:     result.Gofood,
		},
		Maps:      helper.LinkGoogleMaps(result.Latitude, result.Longitude),
		Instagram: helper.LinkInstagram(result.Instagram),
		CreatedBy: &response.UserResponse{
			Name: result.UserCreate.Name,
		},
		UpdatedBy: &response.UserResponse{
			Name: result.UserUpdate.Name,
		},
		Images:    images,
		UpdatedAt: result.UpdatedAt,
	}

	resp.Meta.Status = true
	resp.Meta.Message = "Success fetch coffee shop by id"
	resp.Meta.Errors = nil
	resp.Data = respData

	return c.JSON(resp)
}

// GetCoffeeShops implements [CoffeeShopHandler].
func (f *coffeeShopHandler) GetCoffeeShops(c *fiber.Ctx) error {
	var resp response.DefaultSuccessResponse
	var errResp response.DefaultErrorResponse

	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		errResp.Meta.Status = false
		errResp.Meta.Message = "Unauthorized Access"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusUnauthorized).JSON(errResp)
	}

	results, err := f.CoffeeShopService.GetCoffeeShops(c.Context())
	if err != nil {
		code := "[HANDLER] GetCoffeeShops - 1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Internal server error"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	respDatas := []response.CoffeeShopsResponse{}
	for _, res := range results {
		images := []response.ImagesCoffeeShopResponse{}
		for _, image := range res.Image {
			imag := response.ImagesCoffeeShopResponse{
				Image:     image.Image,
				IsPrimary: image.IsPrimary,
			}
			images = append(images, imag)
		}

		respData := response.CoffeeShopsResponse{
			ID:        res.ID,
			Name:      res.Name,
			Address:   res.Address,
			OpenClose: helper.GenerateOpenTime(res.OpenTime, res.CloseTime),
			Category:  res.Category.Name,
			Images:    images,
		}
		respDatas = append(respDatas, respData)
	}

	resp.Meta.Status = true
	resp.Meta.Message = "Success fetch coffee shops"
	resp.Meta.Errors = nil
	resp.Data = respDatas

	return c.JSON(resp)
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
