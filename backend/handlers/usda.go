package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jwd0526/nutrirecipe/models"
	"github.com/jwd0526/nutrirecipe/services"
)

type USDAHandler struct {
	svc *services.USDAService
}

func NewUSDAHandler(svc *services.USDAService) *USDAHandler {
	return &USDAHandler{svc: svc}
}

type validateRequest struct {
	Ingredients []models.ParsedIngredient `json:"ingredients"`
}

func (h *USDAHandler) Validate(c *gin.Context) {
	var req validateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results := make([]models.ValidatedIngredient, 0, len(req.Ingredients))
	for _, ing := range req.Ingredients {
		results = append(results, h.svc.Validate(c.Request.Context(), ing))
	}
	c.JSON(http.StatusOK, results)
}
