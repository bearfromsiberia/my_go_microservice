package database

import (
	"fmt"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(storagePath string) (*gorm.DB, error) {
	const op = "storage.sqlite.New"

	db, err := gorm.Open(sqlite.Open(storagePath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Отключаем логи
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.AutoMigrate(&Product{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}

type Product struct {
	ID           int       `json:"ID" gorm:"primaryKey"`
	Product_name string    `json:"product_name"`
	Cost         string    `json:"cost"`
	CreatedAt    time.Time `json:"CreatedAt"` // Используйте time.Time
}
