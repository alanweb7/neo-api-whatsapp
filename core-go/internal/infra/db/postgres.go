package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(databaseURL string, isProd bool) (*gorm.DB, error) {
	gLogger := logger.Default
	if isProd {
		gLogger = logger.Default.LogMode(logger.Warn)
	} else {
		gLogger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{Logger: gLogger})
	if err != nil {
		return nil, fmt.Errorf("failed to connect db: %w", err)
	}
	return db, nil
}
