// # SINGLE REASON: Define catalog application input models.
package application

type CreateCategoryInput struct {
	Name string
}

type CreateProductInput struct {
	CategoryID uint
	Name       string
	Attribute  string
	PackSize   float64
}
