// # SINGLE REASON: Define catalog application service contract.
package application

import (
	"context"

	"AnbariAPI/internal/catalog/domain"
)

type Service interface {
	CreateCategory(ctx context.Context, input CreateCategoryInput) (*domain.Category, error)
	GetCategory(ctx context.Context, id uint) (*domain.Category, error)
	ListCategories(ctx context.Context) ([]domain.Category, int64, error)
	DeleteCategory(ctx context.Context, id uint) error
	CreateProduct(ctx context.Context, input CreateProductInput) (*domain.Product, error)
	GetProduct(ctx context.Context, id uint) (*domain.Product, error)
	ListProducts(ctx context.Context) ([]domain.Product, int64, error)
	DeleteProduct(ctx context.Context, id uint) error
}

type CatalogService struct {
	repo domain.Repository
}

func NewCatalogService(repo domain.Repository) *CatalogService {
	return &CatalogService{repo: repo}
}
