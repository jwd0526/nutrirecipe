package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jwd0526/nutrirecipe/models"
)

func mockOllama(t *testing.T, resp models.AgentParseResponse) (*httptest.Server, *AgentService) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content, _ := json.Marshal(resp)
		json.NewEncoder(w).Encode(ollamaResponse{
			Message: ollamaMessage{Role: "assistant", Content: string(content)},
		})
	}))
	return srv, NewAgentService(srv.URL, "test-model")
}

func TestParse_Resolved(t *testing.T) {
	want := models.AgentParseResponse{
		Status: "resolved",
		Ingredients: []models.ParsedIngredient{
			{Name: "chicken breast", QuantityG: 140, Confidence: "high"},
			{Name: "olive oil", QuantityG: 15, Confidence: "high"},
		},
	}
	srv, svc := mockOllama(t, want)
	defer srv.Close()

	resp, err := svc.Parse(models.AgentParseRequest{Input: "2 cups chicken breast\n1 tbsp olive oil"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "resolved" {
		t.Errorf("expected resolved, got %s", resp.Status)
	}
	if len(resp.Ingredients) != 2 {
		t.Errorf("expected 2 ingredients, got %d", len(resp.Ingredients))
	}
}

func TestParse_NeedsClarification(t *testing.T) {
	want := models.AgentParseResponse{
		Status: "needs_clarification",
		Questions: []models.ClarificationQ{{
			ID: "q1", IngredientRaw: "syrup",
			Question: "What type of syrup?",
			Options:  []string{"maple syrup", "corn syrup"},
		}},
	}
	srv, svc := mockOllama(t, want)
	defer srv.Close()

	resp, err := svc.Parse(models.AgentParseRequest{Input: "1 cup syrup"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "needs_clarification" {
		t.Errorf("expected needs_clarification, got %s", resp.Status)
	}
	if len(resp.Questions) == 0 {
		t.Error("expected at least one question")
	}
}

func TestParse_EmptyInput(t *testing.T) {
	want := models.AgentParseResponse{Status: "resolved", Ingredients: []models.ParsedIngredient{}}
	srv, svc := mockOllama(t, want)
	defer srv.Close()

	resp, err := svc.Parse(models.AgentParseRequest{Input: ""})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "resolved" {
		t.Errorf("expected resolved, got %s", resp.Status)
	}
	if len(resp.Ingredients) != 0 {
		t.Errorf("expected 0 ingredients, got %d", len(resp.Ingredients))
	}
}

func TestParse_BadAgentJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(ollamaResponse{
			Message: ollamaMessage{Role: "assistant", Content: "not valid json"},
		})
	}))
	defer srv.Close()
	svc := NewAgentService(srv.URL, "test-model")

	_, err := svc.Parse(models.AgentParseRequest{Input: "chicken breast"})
	if err == nil {
		t.Error("expected error for bad agent JSON response")
	}
}

func TestEvalMatch_Matched(t *testing.T) {
	if evalMatch("chicken breast raw", "chicken breast skinless raw") != "matched" {
		t.Error("expected matched")
	}
}

func TestEvalMatch_LowConfidence(t *testing.T) {
	if evalMatch("chicken breast", "beef sirloin cooked") != "low_confidence" {
		t.Error("expected low_confidence")
	}
}

func TestEvalMatch_EmptyQuery(t *testing.T) {
	if evalMatch("", "chicken breast") != "low_confidence" {
		t.Error("expected low_confidence for empty query")
	}
}

func TestEvalMatch_EmptyResult(t *testing.T) {
	if evalMatch("chicken breast", "") != "low_confidence" {
		t.Error("expected low_confidence for empty result")
	}
}

func TestEvalMatch_ExactMatch(t *testing.T) {
	if evalMatch("whole milk", "whole milk") != "matched" {
		t.Error("expected matched for identical strings")
	}
}
