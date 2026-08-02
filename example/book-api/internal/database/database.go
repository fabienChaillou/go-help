package database

import (
	"github.com/fabienChaillou/book-api/internal/config"
	"github.com/fabienChaillou/book-api/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// New opens a GORM connection and runs automigrations.
func New(cfg config.DBConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&model.Book{}); err != nil {
		return nil, err
	}

	return db, nil
}
