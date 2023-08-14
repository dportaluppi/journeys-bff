package product

import (
	"github.com/dportaluppi/journeys-bff/pkg/recommendations"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type handler struct {
	searcher recommendations.Searcher
}

func NewHandler(productService recommendations.Searcher) *handler {
	return &handler{searcher: productService}
}

func (h *handler) SearchProducts(c *gin.Context) {
	ctx := c.Request.Context()

	query := c.Query("searchTerm")
	storefront := c.Query("storefront")

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	pageNumber, err := strconv.Atoi(c.DefaultQuery("pageNumber", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	f := &recommendations.Filter{
		Storefront: storefront,
		Name:       query,
	}
	result, err := h.searcher.Search(ctx, f, pageSize, pageNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
