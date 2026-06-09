// # SINGLE REASON: Define catalog repository port.
package domain

import "context"

type Repository interface {
	CreateCategory(ctx context.Context, category *Category) error
	GetCategoryByID(ctx context.Context, id uint, preloadProducts bool) (*Category, error)
	ListCategories(ctx context.Context, preloadProducts bool) ([]Category, int64, error)
	DeleteCategory(ctx context.Context, id uint) (bool, error)
	CreateProduct(ctx context.Context, product *Product) error
	GetProductByID(ctx context.Context, id uint, preloadCategory bool) (*Product, error)
	ListProducts(ctx context.Context, preloadCategory bool) ([]Product, int64, error)
	DeleteProduct(ctx context.Context, id uint) (bool, error)
}
