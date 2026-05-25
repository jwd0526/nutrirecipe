package services

import (
	"strings"

	"github.com/jwd0526/nutrirecipe/models"
)

type AgentService struct{}

func NewAgentService() *AgentService {
	return &AgentService{}
}

// Parse implements Agent 1: clarify ambiguous syrup, else resolve with placeholders.
func (a *AgentService) Parse(req models.AgentParseRequest) models.AgentParseResponse {
	if containsSyrupWithoutQualifier(req.Input) {
		return models.AgentParseResponse{
			Status: "needs_clarification",
			Questions: []models.ClarificationQ{
				{
					ID:            "syrup_type",
					IngredientRaw: "syrup",
					Question:      "What type of syrup?",
					Options:       []string{"maple syrup", "corn syrup", "agave syrup", "simple syrup"},
				},
			},
		}
	}
	return models.AgentParseResponse{
		Status:      "resolved",
		Ingredients: parsePlaceholders(req.Input),
	}
}

func containsSyrupWithoutQualifier(input string) bool {
	lower := strings.ToLower(input)
	if !strings.Contains(lower, "syrup") {
		return false
	}
	qualifiers := []string{"maple", "corn", "agave", "simple", "golden", "rice"}
	for _, q := range qualifiers {
		if strings.Contains(lower, q+" syrup") {
			return false
		}
	}
	return true
}

func parsePlaceholders(input string) []models.ParsedIngredient {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	out := make([]models.ParsedIngredient, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		out = append(out, models.ParsedIngredient{
			Name:      line,
			Portion:   1,
			Unit:      "serving",
			QuantityG: 100,
			SearchQueries: models.SearchQueries{
				Primary:      line,
				Alternatives: []string{},
			},
			Confidence: "high",
		})
	}
	return out
}

// evalMatch implements Agent 2: < 2 shared words → low_confidence.
func evalMatch(query, result string) string {
	queryWords := strings.Fields(strings.ToLower(query))
	resultWords := strings.Fields(strings.ToLower(result))

	shared := 0
	for _, qw := range queryWords {
		for _, rw := range resultWords {
			if qw == rw {
				shared++
				break
			}
		}
	}
	if shared < 2 {
		return "low_confidence"
	}
	return "matched"
}
