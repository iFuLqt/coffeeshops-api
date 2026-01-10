package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/ifulqt/coffeeshops-api/config"
	"github.com/ifulqt/coffeeshops-api/internal/adapter/cloudflare"
	"github.com/ifulqt/coffeeshops-api/internal/adapter/handler"
	"github.com/ifulqt/coffeeshops-api/internal/adapter/repository"
	"github.com/ifulqt/coffeeshops-api/internal/core/service"
	"github.com/ifulqt/coffeeshops-api/library/auth"
	"github.com/ifulqt/coffeeshops-api/library/middleware"
)

func RunServer() {
	cfg := config.NewConfig()
	db, err := cfg.ConnectionPostgres()
	if err != nil {
		log.Fatal("Error Connecting Database : %v", err)
		return
	}

	// SETTINGS
	cdfR2 := cfg.LoadAwsConfig()
	s3Client := s3.NewFromConfig(cdfR2)
	r2Adapter := cloudflare.NewCloudFlareR2Adapter(s3Client, cfg)

	jwt := auth.NewJwt(cfg)
	middlewareAuth := middleware.NewMiddleware(jwt)

	// REPOSITORY
	authRepo := repository.NewAuthRepository(db.DB)
	userRepo := repository.NewUserRepository(db.DB)
	categoryRepo := repository.NewCategoryRepository(db.DB)
	coffeeShopRepo := repository.NewCoffeeShopRepository(db.DB)
	uploadImageRepo := repository.NewUploadImageRepository(db.DB)
	facilityRepo := repository.NewFacilityRepository(db.DB)

	// SERVICE
	authServ := service.NewAuthService(authRepo, cfg, jwt)
	userServ := service.NewUserService(userRepo)
	categoryServ := service.NewCategoryService(categoryRepo)
	coffeeShopServ := service.NewCoffeeShopService(coffeeShopRepo)
	uploadImageServ := service.NewUploadImageService(uploadImageRepo, r2Adapter)
	facilityServ := service.NewFacilityService(facilityRepo)

	// HANDLER
	authHandler := handler.NewAuthHandler(authServ)
	userHandler := handler.NewUserHandler(userServ)
	categoryHandler := handler.NewCategoryHandler(categoryServ)
	coffeeShopHandler := handler.NewCoffeeShopHandler(coffeeShopServ)
	uploadImageHandler := handler.NewUploadImageHandler(uploadImageServ)
	facilityHandler := handler.NewFacilityHandler(facilityServ)

	app := fiber.New()
	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${ip} ${status} - ${latency} ${method} ${path}\n",
	}))

	api := app.Group("/api")
	api.Post("/login", authHandler.Login)

	admin := api.Group("/admin")
	admin.Use(middlewareAuth.CheckToken())
	admin.Post("/update-password", userHandler.UpdatePassword)

	category := admin.Group("/categories")
	category.Post("/", categoryHandler.CreateCategory)
	category.Get("/", categoryHandler.GetCategories)
	category.Get("/:categoryID", categoryHandler.GetCategoryByID)
	category.Delete("/:categoryID", categoryHandler.DeleteCategory)
	category.Put("/:categoryID", categoryHandler.UpdateCategory)

	coffeShop := admin.Group("/coffeeshops")
	coffeShop.Post("/", coffeeShopHandler.CreateCoffeeShop)
	coffeShop.Get("/", coffeeShopHandler.GetCoffeeShops)
	coffeShop.Get("/:coffeeshopID", coffeeShopHandler.GetCoffeeShopByID)
	coffeShop.Post("/:coffeeshopID/upload-image", uploadImageHandler.UploadImages)
	coffeShop.Post("/:coffeeshopID/facilities", facilityHandler.CreateFacilityCoffeeShop)

	facility := admin.Group("/facilities")
	facility.Post("/", facilityHandler.CreateFacility)

	go func() {
		if cfg.App.AppPort == "" {
			cfg.App.AppPort = os.Getenv("APP_PORT")
		}
		err := app.Listen(":" + cfg.App.AppPort)
		if err != nil {
			log.Fatal("Error Starting Server : %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	<-quit

	log.Println("Server Shutdown of 5 Seconds")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	app.ShutdownWithContext(ctx)
}
