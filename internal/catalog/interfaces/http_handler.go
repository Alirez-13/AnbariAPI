// # SINGLE REASON: Hold catalog HTTP handler dependencies.
package interfaces

import "AnbariAPI/internal/catalog/application"

type CatalogHandler struct {
	service application.Service
}

func NewCatalogHandler(service application.Service) *CatalogHandler {
	return &CatalogHandler{service: service}
}
