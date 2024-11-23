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

	// Открытие базы данных SQLite
	db, err := gorm.Open(sqlite.Open(storagePath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Отключаем логи
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Миграция схемы для создания таблицы
	err = db.AutoMigrate(&User{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}

// Структура для базы данных
type User struct {
	ID        int       `json:"ID" gorm:"primaryKey"`
	Login     string    `json:"login"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"CreatedAt"` // Используйте time.Time
}
