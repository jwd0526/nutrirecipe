package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jwd0526/nutrirecipe/models"
)

type mockAgentSvc struct{}

func (m *mockAgentSvc) Parse(req models.AgentParseRequest) (models.AgentParseResponse, error) {
	if req.Input == "" {
		return models.AgentParseResponse{Status: "resolved", Ingredients: []models.ParsedIngredient{}}, nil
	}
	lower := strings.ToLower(req.Input)
	if strings.Contains(lower, "syrup") && !strings.Contains(lower, "maple") &&
		!strings.Contains(lower, "corn") && !strings.Contains(lower, "agave") {
		return models.AgentParseResponse{
			Status: "needs_clarification",
			Questions: []models.ClarificationQ{{
				ID: "q1", IngredientRaw: "syrup",
				Question: "What type of syrup?",
				Options:  []string{"maple syrup", "corn syrup", "agave syrup", "simple syrup"},
			}},
		}, nil
	}
	lines := strings.Split(strings.TrimSpace(req.Input), "\n")
	ings := make([]models.ParsedIngredient, 0, len(lines))
	for _, l := range lines {
		if l = strings.TrimSpace(l); l != "" {
			ings = append(ings, models.ParsedIngredient{Name: l, QuantityG: 100, Confidence: "high"})
		}
	}
	return models.AgentParseResponse{Status: "resolved", Ingredients: ings}, nil
}

func init() { gin.SetMode(gin.TestMode) }

func agentRouter() *gin.Engine {
	r := gin.New()
	h := NewAgentHandler(&mockAgentSvc{})
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
