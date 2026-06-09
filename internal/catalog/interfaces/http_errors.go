// # SINGLE REASON: Map catalog errors to HTTP responses.
package interfaces

import (
	"errors"
	"net/http"

	"AnbariAPI/internal/catalog/application"
	"AnbariAPI/internal/catalog/domain"
	"AnbariAPI/internal/catalog/dto"

	"github.com/gin-gonic/gin"
)

func writeCatalogError(c *gin.Context, err error, validationMessage string) {
	if errors.Is(err, domain.ErrNotFound) {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not_found", Message: validationMessage})
		return
	}
	if errors.Is(err, application.ErrValidation) {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "validation_error", Message: err.Error()})
		return
	}
	c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "database_error", Message: validationMessage})
}
