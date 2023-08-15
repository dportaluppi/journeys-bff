package product

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type handler struct {
	searcher Searcher
}

func NewHandler(productService Searcher) *handler {
	return &handler{searcher: productService}
}

func (h *handler) SearchProducts(c *gin.Context) {
	ctx := c.Request.Context()

	f := &Filter{
		Storefront: c.Query("storefront"),
		Name:       c.Query("name"),
		Category:   c.Query("category"),
		SKU:        c.Query("sku"),
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pageNumber, err := strconv.Atoi(c.DefaultQuery("pageNumber", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.searcher.Search(ctx, f, pageSize, pageNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
