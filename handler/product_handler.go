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

func CreateProduct(c *gin.Context) {
	var req dto.ProductCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	db := database.GetDB()
	var category model.Category
	if err := db.First(&category, req.CategoryID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "Category not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to validate category",
		})
		return
	}

	product := model.Product{
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
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to create product",
		})
		return
	}

	db.Preload("Category").First(&product, product.ID)

	c.JSON(http.StatusCreated, dto.ProductResponse{
		ID:           product.ID,
		CategoryID:   product.CategoryID,
		Category:     dto.CategoryResponse{ID: product.Category.ID, Name: product.Category.Name},
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
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid product ID",
		})
		return
	}

	db := database.GetDB()
	var product model.Product
	if err := db.Preload("Category").First(&product, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "Product not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to fetch product",
		})
		return
	}

	c.JSON(http.StatusOK, dto.ProductResponse{
		ID:           product.ID,
		CategoryID:   product.CategoryID,
		Category:     dto.CategoryResponse{ID: product.Category.ID, Name: product.Category.Name},
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
	var products []model.Product
	var total int64

	if err := db.Model(&model.Product{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to fetch products",
		})
		return
	}

	if err := db.Preload("Category").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to fetch products",
		})
		return
	}

	responses := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		responses[i] = dto.ProductResponse{
			ID:           p.ID,
			CategoryID:   p.CategoryID,
			Category:     dto.CategoryResponse{ID: p.Category.ID, Name: p.Category.Name},
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

	c.JSON(http.StatusOK, dto.ProductListResponse{
		Products: responses,
		Total:    total,
	})
}

func DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid product ID",
		})
		return
	}

	db := database.GetDB()
	result := db.Delete(&model.Product{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to delete product",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Product not found",
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Product deleted successfully",
	})
}