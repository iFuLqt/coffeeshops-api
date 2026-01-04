package seeds

import (
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/model"
	"github.com/ifulqt/coffeeshops-api/library/helper"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func UserSeeder(db *gorm.DB) {
	pass, err := helper.HashPassword("admin123")
	if err != nil {
		log.Fatal().Err(err).Msg("Error hashing password")
	}
	admin := model.User{
		Name: "Saiful Anwar",
		Email: "admin@gmail.com",
		Password: pass,
		Role: "admin",
	}

	if err := db.FirstOrCreate(&admin, model.User{Email: "admin@gmail.com"}).Error; err != nil {
		log.Fatal().Err(err).Msg("Error seeding role admin")
	} else {
		log.Info().Msg("Admin role seeding successfully")
	}
}