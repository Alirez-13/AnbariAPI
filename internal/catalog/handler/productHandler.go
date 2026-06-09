// # SINGLE REASON: Preserve product route handler function names.
package handler

import "github.com/gin-gonic/gin"

func CreateProduct(c *gin.Context) {
	defaultHandler().CreateProduct(c)
}

func GetProduct(c *gin.Context) {
	defaultHandler().GetProduct(c)
}

func ListProducts(c *gin.Context) {
	defaultHandler().ListProducts(c)
}

func DeleteProduct(c *gin.Context) {
	defaultHandler().DeleteProduct(c)
}
