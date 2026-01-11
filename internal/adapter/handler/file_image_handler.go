package handler

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/ifulqt/coffeeshops-api/internal/adapter/handler/request"
	"github.com/ifulqt/coffeeshops-api/internal/adapter/handler/response"
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/entity"
	"github.com/ifulqt/coffeeshops-api/internal/core/service"
	"github.com/ifulqt/coffeeshops-api/library/helper"
	"github.com/ifulqt/coffeeshops-api/library/validat"
)

type ImageHandler interface {
	UploadImages(c *fiber.Ctx) error
	DeleteImagesForCoffeeShop(c *fiber.Ctx) error
}

type imageHandler struct {
	ImageService service.ImageService
}

// DeleteImages implements [ImageHandler].
func (u *imageHandler) DeleteImagesForCoffeeShop(c *fiber.Ctx) error {
	var errResp response.DefaultErrorResponse
	var resp response.DefaultSuccessResponse
	var req request.ImageRequest

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
		code := "[HANDLER] DeleteImagesForCoffeeShop - 1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Coffee shop id must be integer"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	err = c.BodyParser(&req)
	if err != nil {
		code := "[HANDLER] DeleteImagesForCoffeeShop - 2"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Invalid request body"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	err = validat.ValidateStruct(&req)
	if err != nil {
		code := "[HANDLER] DeleteImagesForCoffeeShop - 3"
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

	err = u.ImageService.DeleteImagesForCoffeeShop(c.Context(), req.Image, idCoffeeShop)
	if err != nil {
		code := "[HANDLER] DeleteImagesForCoffeeShop - 4"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Internal server error"
		errResp.Meta.Errors = nil
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	resp.Meta.Status = true
	resp.Meta.Message = "Delete image coffee shop successfully"
	resp.Meta.Errors = nil

	return c.JSON(resp)
}

// UploadImages implements [UploadImageHandler].
func (u *imageHandler) UploadImages(c *fiber.Ctx) error {
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

		imageURL, err := u.ImageService.UploadImageR2(c.Context(), uploadReq)
		if err != nil {
			code := "[HANDLER]"
			log.Errorw(code, err)
			errResp.Meta.Status = false
			errResp.Meta.Message = "Failed to upload image"
			errResp.Meta.Errors = nil
			return c.Status(fiber.StatusInternalServerError).JSON(errResp)
		}

		reqEntity := entity.ImageEntity{
			CoffeeShopID: coffeeShopID,
			Image:        imageURL,
			IsPrimary:    i == 0,
		}

		err = u.ImageService.UploadImages(c.Context(), reqEntity)
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

func NewImageHandler(imageServ service.ImageService) ImageHandler {
	return &imageHandler{
		ImageService: imageServ,
	}
}
