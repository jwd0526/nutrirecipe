package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jwd0526/nutrirecipe/models"
)

type agentParser interface {
	Parse(req models.AgentParseRequest) (models.AgentParseResponse, error)
}

type AgentHandler struct {
	svc agentParser
}

func NewAgentHandler(svc agentParser) *AgentHandler {
	return &AgentHandler{svc: svc}
}

func (h *AgentHandler) Parse(c *gin.Context) {
	var req models.AgentParseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.svc.Parse(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "agent failed"})
		return
	}
	c.JSON(http.StatusOK, resp)
}
