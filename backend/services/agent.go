package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/jwd0526/nutrirecipe/models"
)

type AgentService struct {
	ollamaURL string
	model     string
	client    *http.Client
}

func NewAgentService(ollamaURL, model string) *AgentService {
	return &AgentService{
		ollamaURL: ollamaURL,
		model:     model,
		client:    &http.Client{Timeout: 60 * time.Second},
	}
}

const systemPrompt = `You are a recipe ingredient parser. Parse the given ingredient list and respond with ONLY valid JSON.

If all ingredients are clear, respond:
{
  "status": "resolved",
  "ingredients": [
    {
      "name": "chicken breast",
      "portion": 1,
      "unit": "cup",
      "quantity_g": 140,
      "search_queries": {"primary": "cooked chicken breast", "alternatives": ["chicken breast cooked"]},
      "confidence": "high"
    }
  ]
}

If any ingredient type is ambiguous (e.g. "oil", "flour", or "syrup" with no qualifier), respond:
{
  "status": "needs_clarification",
  "questions": [
    {
      "id": "q1",
      "ingredient_raw": "syrup",
      "question": "What type of syrup?",
      "options": ["maple syrup", "corn syrup", "agave syrup", "simple syrup"]
    }
  ]
}

Gram conversion reference:
- Volumes: 1 cup liquid=240g, 1 cup flour=120g, 1 cup cooked rice=195g, 1 cup cooked chicken=140g, 1 tbsp=14g, 1 tsp=5g, 1 oz=28g, 1 lb=454g
- Whole items: 1 large egg=50g, 1 medium egg=44g, 1 medium carrot=61g, 1 large carrot=72g, 1 medium onion=110g, 1 large onion=150g, 1 medium potato=150g, 1 clove garlic=3g, 1 medium tomato=123g, 1 medium apple=182g, 1 medium banana=118g
- Always provide a non-zero quantity_g. If size is unspecified assume medium.
The "primary" search query must be a plain ingredient name suitable for the USDA food database (no quantities or units).
Output nothing outside the JSON.`

type ollamaRequest struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
	Format   string          `json:"format"`
}

type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaResponse struct {
	Message ollamaMessage `json:"message"`
}

func (a *AgentService) Parse(req models.AgentParseRequest) (models.AgentParseResponse, error) {
	userContent := req.Input
	if len(req.QAHistory) > 0 {
		var sb strings.Builder
		sb.WriteString(req.Input)
		sb.WriteString("\n\nClarifications provided:\n")
		for _, qa := range req.QAHistory {
			fmt.Fprintf(&sb, "Q: %s\nA: %s\n", qa.Question, qa.Answer)
		}
		userContent = sb.String()
	}

	payload, err := json.Marshal(ollamaRequest{
		Model: a.model,
		Messages: []ollamaMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userContent},
		},
		Stream: false,
		Format: "json",
	})
	if err != nil {
		return models.AgentParseResponse{}, fmt.Errorf("marshal: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, a.ollamaURL+"/api/chat", bytes.NewReader(payload))
	if err != nil {
		return models.AgentParseResponse{}, fmt.Errorf("build request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(httpReq)
	if err != nil {
		return models.AgentParseResponse{}, fmt.Errorf("ollama request: %w", err)
	}
	defer resp.Body.Close()

	var olResp ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&olResp); err != nil {
		return models.AgentParseResponse{}, fmt.Errorf("decode response: %w", err)
	}

	content := extractJSON(olResp.Message.Content)
	var result models.AgentParseResponse
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return models.AgentParseResponse{}, fmt.Errorf("parse agent JSON: %w", err)
	}
	return result, nil
}

var jsonBlockRe = regexp.MustCompile("(?s)```(?:json)?\\s*(\\{.*?\\})\\s*```")

func extractJSON(s string) string {
	if m := jsonBlockRe.FindStringSubmatch(s); len(m) > 1 {
		return m[1]
	}
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start >= 0 && end > start {
		return s[start : end+1]
	}
	return s
}

// evalMatch is used by usda.go: < 2 shared words → low_confidence.
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
