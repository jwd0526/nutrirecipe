package services

import (
	"strings"
	"testing"
)

func TestMatchStatus_Matched(t *testing.T) {
	s := &USDAService{}
	status := s.matchStatus("chicken breast raw", "chicken breast skinless raw")
	if status != "matched" {
		t.Errorf("expected matched, got %s", status)
	}
}

func TestMatchStatus_LowConfidence(t *testing.T) {
	s := &USDAService{}
	status := s.matchStatus("chicken breast", "beef sirloin cooked")
	if status != "low_confidence" {
		t.Errorf("expected low_confidence, got %s", status)
	}
}

func TestMatchStatus_EmptyQuery(t *testing.T) {
	s := &USDAService{}
	status := s.matchStatus("", "chicken breast")
	if status != "low_confidence" {
		t.Errorf("expected low_confidence for empty query, got %s", status)
	}
}

func TestMatchStatus_EmptyResult(t *testing.T) {
	s := &USDAService{}
	status := s.matchStatus("chicken breast", "")
	if status != "low_confidence" {
		t.Errorf("expected low_confidence for empty result, got %s", status)
	}
}

func TestMatchStatus_ExactMatch(t *testing.T) {
	s := &USDAService{}
	status := s.matchStatus("whole milk", "whole milk")
	if status != "matched" {
		t.Errorf("expected matched for identical strings, got %s", status)
	}
}

func TestUSDASearchURL_Format(t *testing.T) {
	if !strings.HasPrefix(usdaSearchURL, "https://") {
		t.Error("usdaSearchURL must use HTTPS")
	}
	if !strings.Contains(usdaSearchURL, "usda.gov") {
		t.Error("usdaSearchURL must point to usda.gov")
	}
}

func TestNutrientNumbers(t *testing.T) {
	food := usdaFood{
		FdcID:       123,
		Description: "test food",
		FoodNutrients: []usdaNutrient{
			{NutrientNumber: "208", Value: 100},
			{NutrientNumber: "203", Value: 10},
			{NutrientNumber: "204", Value: 5},
			{NutrientNumber: "205", Value: 20},
			{NutrientNumber: "999", Value: 999},
		},
	}

	s := &USDAService{}
	// simulate the extraction logic used in searchUSDA
	opt := extractOption(s, food)

	if opt.Calories != 100 {
		t.Errorf("expected calories 100, got %f", opt.Calories)
	}
	if opt.Protein != 10 {
		t.Errorf("expected protein 10, got %f", opt.Protein)
	}
	if opt.Fat != 5 {
		t.Errorf("expected fat 5, got %f", opt.Fat)
	}
	if opt.Carbs != 20 {
		t.Errorf("expected carbs 20, got %f", opt.Carbs)
	}
}

func TestNutrientNumbers_NoNutrients(t *testing.T) {
	food := usdaFood{FdcID: 1, Description: "empty food"}
	s := &USDAService{}
	opt := extractOption(s, food)
	if opt.Calories != 0 || opt.Protein != 0 || opt.Fat != 0 || opt.Carbs != 0 {
		t.Error("expected all zeros for food with no nutrients")
	}
}
