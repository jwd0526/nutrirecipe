package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jwd0526/nutrirecipe/models"
	"github.com/jwd0526/nutrirecipe/services"
)

func usdaRouter() *gin.Engine {
	r := gin.New()
	h := NewUSDAHandler(services.NewUSDAService(nil, "DEMO_KEY"))
	r.POST("/api/usda/validate", h.Validate)
	return r
}

func TestUSDAValidate_BadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/usda/validate", nil)
	req.Header.Set("Content-Type", "application/json")
	usdaRouter().ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUSDAValidate_EmptyIngredients(t *testing.T) {
	body, _ := json.Marshal(validateRequest{Ingredients: []models.ParsedIngredient{}})
	w := httptest.NewRecorder()
	usdaRouter().ServeHTTP(w, newJSONRequest(http.MethodPost, "/api/usda/validate", body))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var results []models.ValidatedIngredient
	json.Unmarshal(w.Body.Bytes(), &results)
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestUSDAValidate_NilIngredients(t *testing.T) {
	body, _ := json.Marshal(map[string]any{"ingredients": nil})
	w := httptest.NewRecorder()
	usdaRouter().ServeHTTP(w, newJSONRequest(http.MethodPost, "/api/usda/validate", body))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
