package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jwd0526/nutrirecipe/models"
	"github.com/jwd0526/nutrirecipe/services"
)

type AgentHandler struct {
	svc *services.AgentService
}

func NewAgentHandler(svc *services.AgentService) *AgentHandler {
	return &AgentHandler{svc: svc}
}

func (h *AgentHandler) Parse(c *gin.Context) {
	var req models.AgentParseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, h.svc.Parse(req))
}
