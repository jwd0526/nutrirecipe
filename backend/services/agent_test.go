package services

import (
	"testing"

	"github.com/jwd0526/nutrirecipe/models"
)

func TestParse_Resolved(t *testing.T) {
	a := NewAgentService()
	resp := a.Parse(models.AgentParseRequest{Input: "2 cups chicken breast\n1 tbsp olive oil"})
	if resp.Status != "resolved" {
		t.Errorf("expected resolved, got %s", resp.Status)
	}
	if len(resp.Ingredients) != 2 {
		t.Errorf("expected 2 ingredients, got %d", len(resp.Ingredients))
	}
}

func TestParse_NeedsClarification_Syrup(t *testing.T) {
	a := NewAgentService()
	resp := a.Parse(models.AgentParseRequest{Input: "1 cup syrup"})
	if resp.Status != "needs_clarification" {
		t.Errorf("expected needs_clarification, got %s", resp.Status)
	}
	if len(resp.Questions) == 0 {
		t.Error("expected at least one clarification question")
	}
}

func TestParse_QualifiedSyrup_Resolved(t *testing.T) {
	a := NewAgentService()
	resp := a.Parse(models.AgentParseRequest{Input: "2 tbsp maple syrup"})
	if resp.Status != "resolved" {
		t.Errorf("expected resolved for qualified syrup, got %s", resp.Status)
	}
}

func TestParse_EmptyInput(t *testing.T) {
	a := NewAgentService()
	resp := a.Parse(models.AgentParseRequest{Input: ""})
	if resp.Status != "resolved" {
		t.Errorf("expected resolved for empty input, got %s", resp.Status)
	}
	if len(resp.Ingredients) != 0 {
		t.Errorf("expected 0 ingredients for empty input, got %d", len(resp.Ingredients))
	}
}

func TestEvalMatch_Matched(t *testing.T) {
	status := evalMatch("chicken breast raw", "chicken breast skinless raw")
	if status != "matched" {
		t.Errorf("expected matched, got %s", status)
	}
}

func TestEvalMatch_LowConfidence(t *testing.T) {
	status := evalMatch("chicken breast", "beef sirloin cooked")
	if status != "low_confidence" {
		t.Errorf("expected low_confidence, got %s", status)
	}
}

func TestEvalMatch_EmptyQuery(t *testing.T) {
	status := evalMatch("", "chicken breast")
	if status != "low_confidence" {
		t.Errorf("expected low_confidence for empty query, got %s", status)
	}
}

func TestEvalMatch_EmptyResult(t *testing.T) {
	status := evalMatch("chicken breast", "")
	if status != "low_confidence" {
		t.Errorf("expected low_confidence for empty result, got %s", status)
	}
}

func TestEvalMatch_ExactMatch(t *testing.T) {
	status := evalMatch("whole milk", "whole milk")
	if status != "matched" {
		t.Errorf("expected matched for identical strings, got %s", status)
	}
}

func TestContainsSyrup_WithoutQualifier(t *testing.T) {
	if !containsSyrupWithoutQualifier("1 cup syrup") {
		t.Error("expected true for bare syrup")
	}
}

func TestContainsSyrup_WithQualifier(t *testing.T) {
	cases := []string{"maple syrup", "corn syrup", "agave syrup", "simple syrup"}
	for _, c := range cases {
		if containsSyrupWithoutQualifier(c) {
			t.Errorf("expected false for qualified syrup: %s", c)
		}
	}
}

func TestContainsSyrup_NoSyrup(t *testing.T) {
	if containsSyrupWithoutQualifier("chicken breast and rice") {
		t.Error("expected false when syrup not present")
	}
}
