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

type CoffeeShopHandler interface {
	CreateCoffeeShop(c *fiber.Ctx) error
	GetCoffeeShops(c *fiber.Ctx) error
	GetCoffeeShopByID(c *fiber.Ctx) error
	UpdateCoffeeShop(c *fiber.Ctx) error
	DeleteCoffeeShop(c *fiber.Ctx) error

	//FE
	GetCoffeeShopsWithQuery(c *fiber.Ctx) error
}

type coffeeShopHandler struct {
	CoffeeShopService service.CoffeeShopService
	ImageService      service.ImageService
}

// GetCoffeeShopsWithQuery implements [CoffeeShopHandler].
func (f *coffeeShopHandler) GetCoffeeShopsWithQuery(c *fiber.Ctx) error {
	var err error
	var limit int64
	var page int64
	var categoryID int64
	var errResp response.DefaultErrorResponse
	var resp response.DefaultSuccessResponse

	page = 1
	if c.Query("page") != "" {
		page, err = helper.StringToInt(c.Query("page"))
		if err != nil {
			code := "[HANDLER] GetCoffeeShopsWithQuery - 1"
			log.Errorw(code, err)
			errResp.Meta.Status = false
			errResp.Meta.Message = "Invalid page number"
			errResp.Meta.Errors = nil
			return c.Status(fiber.StatusBadRequest).JSON(errResp)
		}
	}

	limit = 6
	if c.Query("limit") != "" {
		limit, err = helper.StringToInt(c.Query("limit"))
		if err != nil {
			code := "[HANDLER] GetCoffeeShopsWithQuery - 2"
			log.Errorw(code, err)
			errResp.Meta.Status = false
			errResp.Meta.Message = "Invalid list number"
			errResp.Meta.Errors = nil
			return c.Status(fiber.StatusBadRequest).JSON(errResp)
		}
	}

	orderBy := "created_at"
	if c.Query("orderBy") != "" {
		orderBy = c.Query("orderBy")
	}

	orderType := "desc"
	if c.Query("orderType") != "" {
		orderType = c.Query("orderType")
	}

	search := ""
	if c.Query("search") != "" {
		search = c.Query("search")
	}

	categoryID = 0
	if c.Query("categoryID") != "" {
		categoryID, err = helper.StringToInt(c.Query("categoryID"))
		if err != nil {
			code := "[HANDLER] GetCoffeeShopsWithQuery - 3"
			log.Errorw(code, err)
			errResp.Meta.Status = false
			errResp.Meta.Message = "category ID bla bla bla"
			errResp.Meta.Errors = nil
			return c.Status(fiber.StatusBadRequest).JSON(errResp)
		}
	}

	reqEntity := entity.QueryString{
		Limit:      int64(limit),
		Page:       int64(page),
		OrderBy:    orderBy,
		OrderType:  orderType,
		Search:     search,
		Status:     true,
		CategoryID: int64(categoryID),
	}

	results, totalData, totalPage, err := f.CoffeeShopService.GetCoffeeShops(c.Context(), reqEntity)
	if err != nil {
		code := "[HANDLER] GetCoffeeShopsWithQuery - 4"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Internal server error"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	respDatas := make([]response.CoffeeShopsResponse, len(results))
	for i, res := range results {
		respImages := []response.ImagesCoffeeShopResponse{}
		for _, image := range res.Image {
			respImage := response.ImagesCoffeeShopResponse{
				Image: image.Image,
			}
			respImages = append(respImages, respImage)
		}

		respDatas[i] = response.CoffeeShopsResponse{
			ID:        res.ID,
			Name:      res.Name,
			Slug:      res.Slug,
			Address:   res.Address,
			OpenClose: helper.GenerateOpenTime(res.OpenTime, res.CloseTime),
			Category:  res.Category.Name,
			Images:    respImages,
		}
	}

	resp.Meta.Status = true
	resp.Meta.Message = "Fetched data coffee shops successfully"
	resp.Meta.Errors = nil
	resp.Data = respDatas
	resp.Pagination = &response.PaginationResponse{
		TotalRecords: totalData,
		Page:         int64(page),
		PerPage:      int64(limit),
		TotalPages:   totalPage,
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

	active := true
	reqEntity := entity.CoffeeShopEntity{
		Name:      req.CoffeeShop,
		Address:   req.Address,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		OpenTime:  req.OpenTime,
		CloseTime: req.CloseTime,
		Instagram: req.Instagram,
		UserCreate: entity.UserEntity{
			ID: int64(userID),
		},
		UserUpdate: entity.UserEntity{
			ID: int64(userID),
		},
		Category: entity.CategoryEntity{
			ID: req.CategoryID,
		},
		IsActive: &active,
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
	var errResp response.DefaultErrorResponse
	var resp response.DefaultSuccessResponse

	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		errResp.Meta.Status = false
		errResp.Meta.Message = "Unauthorized access"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusUnauthorized).JSON(errResp)
	}

	idParam := c.Params("coffeeshopID")
	idCoffeeShop, err := helper.StringToInt(idParam)
	if err != nil {
		code := "[HANDLER] DeleteCoffeeShop - 1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Coffee shop ID must be integer"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	err = f.ImageService.DeleteImages(c.Context(), int64(idCoffeeShop))
	if err != nil {
		code := "[HANDLER] DeleteCoffeeShop - 2"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Internal server error"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	err = f.CoffeeShopService.DeleteCoffeeShop(c.Context(), int64(idCoffeeShop))
	if err != nil {
		code := "[HANDLER] DeleteCoffeeShop - 3"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Errors = nil
		if errors.Is(err, domerror.ErrDataNotFound) {
			errResp.Meta.Message = "Coffee shop not found"
			return c.Status(fiber.StatusNotFound).JSON(errResp)
		}
		errResp.Meta.Message = "Internal server error"
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	resp.Meta.Status = true
	resp.Meta.Message = "Delete coffee shop successfully"
	resp.Meta.Errors = nil

	return c.JSON(resp)
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

	result, err := f.CoffeeShopService.GetCoffeeShopByID(c.Context(), int64(id))
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

	facilities := []response.FacilityResponse{}
	for _, fac := range result.Facility {
		facility := response.FacilityResponse{
			Name: fac.Name,
		}
		facilities = append(facilities, facility)
	}

	respData := response.CoffeeShopByIDResponse{
		ID:        result.ID,
		Name:      result.Name,
		Slug:      result.Slug,
		Address:   result.Address,
		OpenClose: helper.GenerateOpenTime(result.OpenTime, result.CloseTime),
		Facility:  facilities,
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

	reqEntity := entity.QueryString{
		Limit:      10,
		Page:       1,
		OrderBy:    "",
		OrderType:  "",
		Search:     "",
		CategoryID: 0,
	}

	results, _, _, err := f.CoffeeShopService.GetCoffeeShops(c.Context(), reqEntity)
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
			Slug:      res.Slug,
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
	var resp response.DefaultSuccessResponse
	var errResp response.DefaultErrorResponse
	var req request.UpdateCoffeeShopRequest

	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		errResp.Meta.Status = false
		errResp.Meta.Message = "Unauthorized access"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusUnauthorized).JSON(errResp)
	}

	idParameter := c.Params("coffeeshopID")
	idCoffeeShop, err := helper.StringToInt(idParameter)
	if err != nil || idCoffeeShop <= 0 {
		code := "[HANDLER] UpdateCoffeeShop - 1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Coffee shop ID must be integer"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	err = c.BodyParser(&req)
	if err != nil {
		code := "[HANDLER] UpdateCoffeeShop - 2"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Invalid request body"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	err = validat.ValidateStruct(&req)
	if err != nil {
		code := "[HANDLER] UpdateCoffeeShop - 3"
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
		Name:      req.CoffeeShop,
		Address:   req.Address,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		OpenTime:  req.OpenTime,
		CloseTime: req.CloseTime,
		Instagram: req.Instagram,
		UserUpdate: entity.UserEntity{
			ID: int64(userID),
		},
		IsActive: req.IsActive,
		Category: entity.CategoryEntity{
			ID: req.CategoryID,
		},
	}

	err = f.CoffeeShopService.UpdateCoffeeShop(c.Context(), reqEntity, idCoffeeShop)
	if err != nil {
		code := "[HANDLER] UpdateCoffeeShop - 4"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Errors = nil
		if errors.Is(err, domerror.ErrDataNotFound) {
			errResp.Meta.Message = "Coffee shop not found"
			return c.Status(fiber.StatusNotFound).JSON(errResp)
		}
		errResp.Meta.Message = "Internal server error"
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	resp.Meta.Status = true
	resp.Meta.Message = "Update coffee shops successfully"
	resp.Meta.Errors = nil

	return c.JSON(resp)
}

func NewCoffeeShopHandler(coffeeShopServ service.CoffeeShopService, imageService service.ImageService) CoffeeShopHandler {
	return &coffeeShopHandler{
		CoffeeShopService: coffeeShopServ,
		ImageService:      imageService,
	}
}
