package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type User struct {
	ID        int       `json:"ID" gorm:"primaryKey"`
	Login     string    `json:"login"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"CreatedAt"` // Используйте time.Time
}

// CreateUser создает нового пользователя
func CreateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		type CreateUserRequest struct {
			Login    string `json:"login" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		var req CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}

		// Создаем нового пользователя с использованием GORM
		user := User{Login: req.Login, Password: req.Password}
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
			return
		}

		// Возвращаем созданного пользователя
		c.JSON(http.StatusCreated, gin.H{"id": user.ID})
	}
}

func GetUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []User

		// Получаем всех пользователей с использованием GORM
		if err := db.Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users: " + err.Error()})
			return
		}

		// Формируем структуру данных для ответа
		var usersResponse []gin.H
		for _, user := range users {
			usersResponse = append(usersResponse, gin.H{
				"ID":        user.ID,
				"Login":     user.Login,
				"Password":  user.Password,
				"CreatedAt": user.CreatedAt.Format("2006-01-02 15:04:05"), // Преобразуем дату в строку
			})
		}

		// Возвращаем список пользователей
		c.JSON(http.StatusOK, gin.H{"users": usersResponse})
	}
}

// UpdateUser обновляет данные пользователя
func UpdateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		type UpdateUserRequest struct {
			Login    string `json:"login"`
			Password string `json:"password" binding:"omitempty"`
		}

		var req UpdateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}

		// Ищем пользователя по ID
		var user User
		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// Обновляем поля пользователя
		if req.Login != "" {
			user.Login = req.Login
		}
		if req.Password != "" {
			user.Password = req.Password
		}

		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
	}
}

// DeleteUser удаляет пользователя
func DeleteUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// Ищем пользователя по ID
		var user User
		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// Удаляем пользователя
		if err := db.Delete(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}
