// # SINGLE REASON: Persist catalog entities with GORM.
package infrastructure

import (
	"context"
	"errors"

	"AnbariAPI/internal/catalog/domain"

	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) CreateCategory(ctx context.Context, category *domain.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *GormRepository) GetCategoryByID(ctx context.Context, id uint, preloadProducts bool) (*domain.Category, error) {
	var category domain.Category
	query := r.db.WithContext(ctx)
	if preloadProducts {
		query = query.Preload("Products")
	}
	if err := query.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &category, nil
}

func (r *GormRepository) ListCategories(ctx context.Context, preloadProducts bool) ([]domain.Category, int64, error) {
	var categories []domain.Category
	var total int64
	db := r.db.WithContext(ctx)
	if err := db.Model(&domain.Category{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	query := db
	if preloadProducts {
		query = query.Preload("Products")
	}
	if err := query.Find(&categories).Error; err != nil {
		return nil, 0, err
	}
	return categories, total, nil
}

func (r *GormRepository) DeleteCategory(ctx context.Context, id uint) (bool, error) {
	result := r.db.WithContext(ctx).Delete(&domain.Category{}, id)
	return result.RowsAffected > 0, result.Error
}

func (r *GormRepository) CreateProduct(ctx context.Context, product *domain.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *GormRepository) GetProductByID(ctx context.Context, id uint, preloadCategory bool) (*domain.Product, error) {
	var product domain.Product
	query := r.db.WithContext(ctx)
	if preloadCategory {
		query = query.Preload("Category")
	}
	if err := query.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *GormRepository) ListProducts(ctx context.Context, preloadCategory bool) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64
	db := r.db.WithContext(ctx)
	if err := db.Model(&domain.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	query := db
	if preloadCategory {
		query = query.Preload("Category")
	}
	if err := query.Find(&products).Error; err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (r *GormRepository) DeleteProduct(ctx context.Context, id uint) (bool, error) {
	result := r.db.WithContext(ctx).Delete(&domain.Product{}, id)
	return result.RowsAffected > 0, result.Error
}
