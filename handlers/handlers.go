package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Product struct {
	ID           int       `json:"ID" gorm:"primaryKey"`
	Product_name string    `json:"product_name"`
	Cost         string    `json:"cost"`
	CreatedAt    time.Time `json:"CreatedAt"`
}

func CreateProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		type CreateProductRequest struct {
			Product_name string `json:"product_name" binding:"required"`
			Cost         string `json:"cost" binding:"required"`
		}

		var req CreateProductRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}

		Product := Product{Product_name: req.Product_name, Cost: req.Cost}
		if err := db.Create(&Product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Product: " + err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"id": Product.ID})
	}
}

func GetProducts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var Products []Product

		if err := db.Find(&Products).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Products: " + err.Error()})
			return
		}

		var ProductsResponse []gin.H
		for _, Product := range Products {
			ProductsResponse = append(ProductsResponse, gin.H{
				"ID":           Product.ID,
				"Product_name": Product.Product_name,
				"Cost":         Product.Cost,
				"CreatedAt":    Product.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}

		c.JSON(http.StatusOK, gin.H{"Products": ProductsResponse})
	}
}

func UpdateProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ProductID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Product ID"})
			return
		}

		type UpdateProductRequest struct {
			Product_name string `json:"product_name"`
			Cost         string `json:"cost" binding:"omitempty"`
		}

		var req UpdateProductRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}

		var Product Product
		if err := db.First(&Product, ProductID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		if req.Product_name != "" {
			Product.Product_name = req.Product_name
		}
		if req.Cost != "" {
			Product.Cost = req.Cost
		}

		if err := db.Save(&Product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Product: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
	}
}

func DeleteProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ProductID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Product ID"})
			return
		}

		var Product Product
		if err := db.First(&Product, ProductID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		if err := db.Delete(&Product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Product: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
	}
}
