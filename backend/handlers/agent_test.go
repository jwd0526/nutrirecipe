package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jwd0526/nutrirecipe/models"
	"github.com/jwd0526/nutrirecipe/services"
)

func init() { gin.SetMode(gin.TestMode) }

func agentRouter() *gin.Engine {
	r := gin.New()
	h := NewAgentHandler(services.NewAgentService())
	r.POST("/api/agent/parse", h.Parse)
	return r
}

func TestAgentParse_Resolved(t *testing.T) {
	body, _ := json.Marshal(models.AgentParseRequest{Input: "chicken breast\nrice"})
	w := httptest.NewRecorder()
	agentRouter().ServeHTTP(w, newJSONRequest(http.MethodPost, "/api/agent/parse", body))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp models.AgentParseResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Status != "resolved" {
		t.Errorf("expected resolved, got %s", resp.Status)
	}
	if len(resp.Ingredients) != 2 {
		t.Errorf("expected 2 ingredients, got %d", len(resp.Ingredients))
	}
}

func TestAgentParse_NeedsClarification(t *testing.T) {
	body, _ := json.Marshal(models.AgentParseRequest{Input: "1 cup syrup"})
	w := httptest.NewRecorder()
	agentRouter().ServeHTTP(w, newJSONRequest(http.MethodPost, "/api/agent/parse", body))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp models.AgentParseResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Status != "needs_clarification" {
		t.Errorf("expected needs_clarification, got %s", resp.Status)
	}
}

func TestAgentParse_BadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/agent/parse", bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")
	agentRouter().ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func newJSONRequest(method, path string, body []byte) *http.Request {
	req, _ := http.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}
