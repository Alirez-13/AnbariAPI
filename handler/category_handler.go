package handler

import (
	"AnbariAPI/database"
	"AnbariAPI/dto"
	"AnbariAPI/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateCategory(c *gin.Context) {
	var req dto.CategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	category := model.Category{
		Name: req.Name,
	}

	db := database.GetDB()
	if err := db.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to create category",
		})
		return
	}

	c.JSON(http.StatusCreated, dto.CategoryResponse{
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
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid category ID",
		})
		return
	}

	db := database.GetDB()
	var category model.Category
	if err := db.Preload("Products").First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "Category not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to fetch category",
		})
		return
	}

	products := make([]dto.ProductResponse, len(category.Products))
	for i, p := range category.Products {
		products[i] = dto.ProductResponse{
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
		"products":  products,
		"created_at": category.CreatedAt,
		"updated_at": category.UpdatedAt,
	})
}

func ListCategories(c *gin.Context) {
	db := database.GetDB()
	var categories []model.Category
	var total int64

	if err := db.Model(&model.Category{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to fetch categories",
		})
		return
	}

	if err := db.Preload("Products").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to fetch categories",
		})
		return
	}

	responses := make([]dto.CategoryResponse, len(categories))
	for i, cat := range categories {
		responses[i] = dto.CategoryResponse{
			ID:        cat.ID,
			Name:      cat.Name,
			CreatedAt: cat.CreatedAt,
			UpdatedAt: cat.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, dto.CategoryListResponse{
		Categories: responses,
		Total:      total,
	})
}

func DeleteCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid category ID",
		})
		return
	}

	db := database.GetDB()
	result := db.Delete(&model.Category{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to delete category",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Category not found",
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Category deleted successfully",
	})
}