package journey

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	CreateJourney(c *gin.Context)
	GetJourneys(c *gin.Context)
	GetJourneyByID(c *gin.Context)
}

type handler struct {
	service Service
}

func NewHandler(s Service) Handler {
	return &handler{service: s}
}

func (h *handler) CreateJourney(c *gin.Context) {
	var j JourneyWriteModel

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

func (h *handler) GetJourneys(c *gin.Context) {
	accountId := c.Param("accountId")

	journeys, err := h.service.GetJourneys(c.Request.Context(), accountId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, journeys)
}

func (h *handler) GetJourneyByID(c *gin.Context) {
	accountId := c.Param("accountId")
	journeyID := c.Param("journeyId")

	journey, err := h.service.GetJourneyByID(c.Request.Context(), accountId, journeyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, journey)
}
