package config

import (
	"sync/atomic"

	"gorm.io/gorm"
)

type Config struct {
	FileserverHits atomic.Int32
	DB             *gorm.DB
	Platform       string
	JWTSecret      string
}

func New(db *gorm.DB, platform string, jwtSecret string) *Config {
	return &Config{
		DB:       db,
		Platform: platform,
		JWTSecret: jwtSecret,
	}
}
