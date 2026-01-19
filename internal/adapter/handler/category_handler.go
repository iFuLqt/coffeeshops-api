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

type CategoryHandler interface {
	GetCategories(c *fiber.Ctx) error
	GetCategoryByID(c *fiber.Ctx) error
	CreateCategory(c *fiber.Ctx) error
	UpdateCategory(c *fiber.Ctx) error
	DeleteCategory(c *fiber.Ctx) error
}

type categoryHandler struct {
	CategoryService service.CategoryService
}

// CreateCategory implements [CategoryHandler].
func (f *categoryHandler) CreateCategory(c *fiber.Ctx) error {
	var err error
	var errResp response.DefaultErrorResponse
	var resp response.DefaultSuccessResponse
	var req request.CategoryRequest

	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		errResp.Meta.Status = false
		errResp.Meta.Message = "Unauthorized access"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusUnauthorized).JSON(errResp)
	}

	err = c.BodyParser(&req)
	if err != nil {
		code := "[HANDLER] CreateCategory - 1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Invalid request body"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	err = validat.ValidateStruct(&req)
	if err != nil {
		code := "[HANDLER] CreateCategory - 2"
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

	reqEntity := entity.CategoryEntity{
		Name: req.Category,
		CreatedBy: entity.UserEntity{
			ID: int64(userID),
		},
	}

	err = f.CategoryService.CreateCategory(c.Context(), reqEntity)
	if err != nil {
		code := "[HANDLER] CreateCategory - 3"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Internal server error"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	resp.Meta.Status = true
	resp.Meta.Message = "Create category successfully"
	resp.Meta.Errors = nil

	return c.JSON(resp)
}

// DeleteCategory implements [CategoryHandler].
func (f *categoryHandler) DeleteCategory(c *fiber.Ctx) error {
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

	idParameter := c.Params("categoryID")
	id, err := helper.StringToInt(idParameter)
	if err != nil || id <= 0 {
		code := "[HANDLER] DeleteCategory - 1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Category ID must be an integer"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	err = f.CategoryService.DeleteCategory(c.Context(), int64(id))
	if err != nil {
		code := "[HANDLER] DeleteCategory - 2"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Errors = nil
		if errors.Is(err, domerror.ErrDeleteCategory) {
			errResp.Meta.Message = "Category is still used by coffee shops"
			return c.Status(fiber.StatusBadRequest).JSON(errResp)
		} else if errors.Is(err, domerror.ErrDataNotFound) {
			errResp.Meta.Message = "Category not found"
			return c.Status(fiber.StatusNotFound).JSON(errResp)
		}
		errResp.Meta.Message = "Internal server error"
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	resp.Meta.Status = true
	resp.Meta.Message = "Delete category successfully"
	resp.Meta.Errors = nil

	return c.JSON(resp)
}

// GetCategories implements [CategoryHandler].
func (f *categoryHandler) GetCategories(c *fiber.Ctx) error {
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

	results, err := f.CategoryService.GetCategories(c.Context())
	if err != nil {
		code := "[HANDLER] GetCategories - 1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Data not found"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	categoryResps := []response.CategoryResponse{}
	for _, val := range results {
		categoryResp := response.CategoryResponse{
			ID:       val.ID,
			Category: val.Name,
			Slug:     val.Slug,
			CreatedBy: response.UserResponse{
				ID:   val.CreatedBy.ID,
				Name: val.CreatedBy.Name,
			},
		}
		categoryResps = append(categoryResps, categoryResp)
	}

	resp.Meta.Status = true
	resp.Meta.Message = "Categories fetched successfully"
	resp.Meta.Errors = nil
	resp.Data = categoryResps

	return c.JSON(resp)
}

// GetCategoryByID implements [CategoryHandler].
func (f *categoryHandler) GetCategoryByID(c *fiber.Ctx) error {
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

	idParameter := c.Params("categoryID")
	id, err := helper.StringToInt(idParameter)
	if err != nil || id <= 0 {
		code := "[HANDLER] GetCategoryByID - 1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Category ID must be an integer"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	result, err := f.CategoryService.GetCategoryByID(c.Context(), int64(id))
	if err != nil {
		code := "[HANDLER] GetCategoryByID - 2"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Internal server error"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	categoryResp := response.CategoryResponse{
		ID:       result.ID,
		Category: result.Name,
		Slug:     result.Slug,
		CreatedBy: response.UserResponse{
			ID:   result.CreatedBy.ID,
			Name: result.CreatedBy.Name,
		},
	}

	resp.Meta.Status = true
	resp.Meta.Message = "Category fetched successfully"
	resp.Meta.Errors = nil
	resp.Data = categoryResp

	return c.JSON(resp)
}

// UpdateCategory implements [CategoryHandler].
func (f *categoryHandler) UpdateCategory(c *fiber.Ctx) error {
	var req request.CategoryRequest
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

	idParameter := c.Params("categoryID")
	id, err := helper.StringToInt(idParameter)
	if err != nil || id <= 0{
		code := "[HANDLER] UpdateCategory - 1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Category ID must be an integer"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	err = c.BodyParser(&req)
	if err != nil {
		code := "[HANDLER] UpdateCategory - 2"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Invalid request body"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	err = validat.ValidateStruct(&req)
	if err != nil {
		code := "[HANDLER] UpdateCategory - 3"
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

	reqEntity := entity.CategoryEntity{
		ID:   int64(id),
		Name: req.Category,
	}

	err = f.CategoryService.UpdateCategory(c.Context(), reqEntity)
	if err != nil {
		code := "[HANDLER] UpdateCategory - 4"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Errors = nil
		if errors.Is(err, domerror.ErrDataNotFound) {
			errResp.Meta.Message = "Category not found"
			return c.Status(fiber.StatusBadRequest).JSON(errResp)
		}
		errResp.Meta.Message = "Internal server error"
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	resp.Meta.Status = true
	resp.Meta.Message = "Update category successfully"
	resp.Meta.Errors = nil

	return c.JSON(resp)
}

func NewCategoryHandler(categoryServ service.CategoryService) CategoryHandler {
	return &categoryHandler{
		CategoryService: categoryServ,
	}
}
