package seeds

import (
	"github.com/ifulqt/coffeeshops-api/internal/core/domain/model"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func FacilitySeed(db *gorm.DB) {
	facilities := []model.Facility{
		{
			Name: "WiFi",
			Code: "wifi",
		},
		{
			Name: "Area Parkir",
			Code: "parking",
		},
		{
			Name: "Musholla",
			Code: "prayer_room",
		},
		{
			Name: "Ruang Merokok",
			Code: "smoking_room",
		},
		{
			Name: "Toilet",
			Code: "toilet",
		},
		{
			Name: "Area Outdoor",
			Code: "outdoor_area",
		},
		{
			Name: "Area Indoor",
			Code: "indoor_area",
		},
	}

	for _, facility := range facilities {
		err := db.Where("code = ?", facility.Code).FirstOrCreate(&facility).Error
		if err != nil {
			log.Fatal().Err(err).Msg("Error seeding facility")
		} else {
			log.Info().Msg("Seeding facility successfully")
		}
	}
}
