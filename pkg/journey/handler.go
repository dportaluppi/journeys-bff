package journey

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	CreateJourney(c *gin.Context)
}

type handler struct {
	service Service
}

func NewHandler(s Service) Handler {
	return &handler{service: s}
}

func (h *handler) CreateJourney(c *gin.Context) {
	var j Journey

	if err := c.ShouldBindJSON(&j); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	j.AccountID = c.Param("accountId")
	newJourney, err := h.service.CreateJourney(c.Request.Context(), &j)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newJourney)
}
