package segment

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type handler struct {
	getter Getter
}

func NewHandler(audienceService Getter) *handler {
	return &handler{getter: audienceService}
}

func (h *handler) GetSegments(c *gin.Context) {
	ctx := c.Request.Context()

	provider := c.Query("provider")
	name := c.Query("name")
	id := c.Query("id")

	pageSize, err := strconv.Atoi(c.DefaultQuery("size", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pageNumber, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	f := &Filter{
		Provider: provider,
		Name:     name,
		ID:       id,
	}
	result, err := h.getter.GetSegments(ctx, f, pageSize, pageNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
