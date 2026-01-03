package config

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	DB *gorm.DB
}

func (cfg Config) ConnectionPostgres() (*Postgres, error) {
	dbConnString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.Psql.DBUser,
		cfg.Psql.DBPassword,
		cfg.Psql.DBHost,
		cfg.Psql.DBPort,
		cfg.Psql.DBName)
	db, err := gorm.Open(postgres.Open(dbConnString), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("[ConnectionPostgres - 1] Failed to connect database")
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Error().Err(err).Msg("[ConnectionPostgres - 2] Failed to get database")
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.Psql.DBMaxIdle)
	sqlDB.SetMaxOpenConns(cfg.Psql.DBMaxOpen)
	
	return &Postgres{DB: db}, nil
}
