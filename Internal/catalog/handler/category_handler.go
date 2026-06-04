package handler

import (
	"AnbariAPI/Internal/catalog/domain"
	dto2 "AnbariAPI/Internal/catalog/dto"
	"AnbariAPI/shared/database"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateCategory(c *gin.Context) {
	var req dto2.CategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto2.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	category := domain.Category{
		Name: req.Name,
	}

	db := database.GetDB()
	if err := db.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto2.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to create category",
		})
		return
	}

	c.JSON(http.StatusCreated, dto2.CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	})
}

func GetCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto2.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid category ID",
		})
		return
	}

	db := database.GetDB()
	var category domain.Category
	if err := db.Preload("Products").First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto2.ErrorResponse{
				Error:   "not_found",
				Message: "Category not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto2.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to fetch category",
		})
		return
	}

	products := make([]dto2.ProductResponse, len(category.Products))
	for i, p := range category.Products {
		products[i] = dto2.ProductResponse{
			ID:           p.ID,
			CategoryID:   p.CategoryID,
			Name:         p.Name,
			Attribute:    p.Attribute,
			Unit:         p.Unit,
			PackSize:     p.PackSize,
			IsPackable:   p.IsPackable,
			BaseUnit:     p.BaseUnit,
			CurrentStock: p.CurrentStock,
			DisplayStock: p.DisplayStock,
			DisplayUnit:  p.DisplayUnit,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         category.ID,
		"name":       category.Name,
		"products":   products,
		"created_at": category.CreatedAt,
		"updated_at": category.UpdatedAt,
	})
}

func ListCategories(c *gin.Context) {
	db := database.GetDB()
	var categories []domain.Category
	var total int64

	if err := db.Model(&domain.Category{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto2.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to fetch categories",
		})
		return
	}

	if err := db.Preload("Products").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto2.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to fetch categories",
		})
		return
	}

	responses := make([]dto2.CategoryResponse, len(categories))
	for i, cat := range categories {
		responses[i] = dto2.CategoryResponse{
			ID:        cat.ID,
			Name:      cat.Name,
			CreatedAt: cat.CreatedAt,
			UpdatedAt: cat.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, dto2.CategoryListResponse{
		Categories: responses,
		Total:      total,
	})
}

func DeleteCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto2.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid category ID",
		})
		return
	}

	db := database.GetDB()
	result := db.Delete(&domain.Category{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, dto2.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to delete category",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, dto2.ErrorResponse{
			Error:   "not_found",
			Message: "Category not found",
		})
		return
	}

	c.JSON(http.StatusOK, dto2.SuccessResponse{
		Message: "Category deleted successfully",
	})
}
