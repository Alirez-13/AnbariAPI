package handler

import (
	dto2 "AnbariAPI/catalog/dto"
	"AnbariAPI/shared/database"
	"AnbariAPI/shared/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateProduct(c *gin.Context) {
	var req dto2.ProductCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto2.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	db := database.GetDB()
	var category models.Category
	if err := db.First(&category, req.CategoryID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, dto2.ErrorResponse{
				Error:   "validation_error",
				Message: "Category not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto2.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to validate category",
		})
		return
	}

	product := models.Product{
		CategoryID:   req.CategoryID,
		Name:         req.Name,
		Attribute:    req.Attribute,
		Unit:         req.Unit,
		PackSize:     req.PackSize,
		IsPackable:   req.IsPackable,
		BaseUnit:     req.BaseUnit,
		CurrentStock: req.CurrentStock,
	}

	if err := db.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto2.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to create product",
		})
		return
	}

	db.Preload("Category").First(&product, product.ID)

	c.JSON(http.StatusCreated, dto2.ProductResponse{
		ID:           product.ID,
		CategoryID:   product.CategoryID,
		Category:     dto2.CategoryResponse{ID: product.Category.ID, Name: product.Category.Name},
		Name:         product.Name,
		Attribute:    product.Attribute,
		Unit:         product.Unit,
		PackSize:     product.PackSize,
		IsPackable:   product.IsPackable,
		BaseUnit:     product.BaseUnit,
		CurrentStock: product.CurrentStock,
		DisplayStock: product.DisplayStock,
		DisplayUnit:  product.DisplayUnit,
		CreatedAt:    product.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    product.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

func GetProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto2.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid product ID",
		})
		return
	}

	db := database.GetDB()
	var product models.Product
	if err := db.Preload("Category").First(&product, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto2.ErrorResponse{
				Error:   "not_found",
				Message: "Product not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto2.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to fetch product",
		})
		return
	}

	c.JSON(http.StatusOK, dto2.ProductResponse{
		ID:           product.ID,
		CategoryID:   product.CategoryID,
		Category:     dto2.CategoryResponse{ID: product.Category.ID, Name: product.Category.Name},
		Name:         product.Name,
		Attribute:    product.Attribute,
		Unit:         product.Unit,
		PackSize:     product.PackSize,
		IsPackable:   product.IsPackable,
		BaseUnit:     product.BaseUnit,
		CurrentStock: product.CurrentStock,
		DisplayStock: product.DisplayStock,
		DisplayUnit:  product.DisplayUnit,
		CreatedAt:    product.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    product.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

func ListProducts(c *gin.Context) {
	db := database.GetDB()
	var products []models.Product
	var total int64

	if err := db.Model(&models.Product{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto2.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to fetch products",
		})
		return
	}

	if err := db.Preload("Category").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto2.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to fetch products",
		})
		return
	}

	responses := make([]dto2.ProductResponse, len(products))
	for i, p := range products {
		responses[i] = dto2.ProductResponse{
			ID:           p.ID,
			CategoryID:   p.CategoryID,
			Category:     dto2.CategoryResponse{ID: p.Category.ID, Name: p.Category.Name},
			Name:         p.Name,
			Attribute:    p.Attribute,
			Unit:         p.Unit,
			PackSize:     p.PackSize,
			IsPackable:   p.IsPackable,
			BaseUnit:     p.BaseUnit,
			CurrentStock: p.CurrentStock,
			DisplayStock: p.DisplayStock,
			DisplayUnit:  p.DisplayUnit,
			CreatedAt:    p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:    p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	c.JSON(http.StatusOK, dto2.ProductListResponse{
		Products: responses,
		Total:    total,
	})
}

func DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto2.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid product ID",
		})
		return
	}

	db := database.GetDB()
	result := db.Delete(&models.Product{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, dto2.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to delete product",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, dto2.ErrorResponse{
			Error:   "not_found",
			Message: "Product not found",
		})
		return
	}

	c.JSON(http.StatusOK, dto2.SuccessResponse{
		Message: "Product deleted successfully",
	})
}
