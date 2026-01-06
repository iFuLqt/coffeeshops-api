package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/ifulqt/coffeeshops-api/config"
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
	jwt := auth.NewJwt(cfg)
	middlewareAuth := middleware.NewMiddleware(jwt)

	// REPOSITORY
	authRepo := repository.NewAuthRepository(db.DB)
	userRepo := repository.NewUserRepository(db.DB)

	// SERVICE
	authServ := service.NewAuthService(authRepo, cfg, jwt)
	userServ := service.NewUserService(userRepo)

	// HANDLER
	authHandler := handler.NewAuthHandler(authServ)
	userHandler := handler.NewUserHandler(userServ)

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
