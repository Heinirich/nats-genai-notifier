package ai

import (
	"Moringa_AI/src/supportai/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// ollamaURL is the local Ollama endpoint.
const ollamaURL = "http://localhost:11434/api/generate"

// ollamaRequest is the payload sent to Ollama.
type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// ollamaResponse is the non-streaming response from Ollama.
type ollamaResponse struct {
	Response string `json:"response"`
}

// GenerateSupportResponse sends the raw user message to the local LLM
// and expects JSON describing a Ticket.
func GenerateSupportResponse(raw string) (*models.Ticket, error) {
	prompt := fmt.Sprintf(`
You are an AI support assistant.

Given a raw user complaint, you must:
- Create a clear, concise title.
- Rewrite the body in a professional tone.
- Assign a priority: "High", "Medium", or "Low".
- Suggest a recommended action for the support team.

Return ONLY valid JSON with the following structure (no extra text):

{
  "title": "...",
  "body": "...",
  "priority": "High" | "Medium" | "Low",
  "action": "..."
}

User message:
%s
`, raw)

	reqBody, err := json.Marshal(ollamaRequest{
		Model:  "gemma3", // Name of the model you pulled with Ollama
		Prompt: prompt,
		Stream: false,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal ollama request: %w", err)
	}

	resp, err := http.Post(ollamaURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("call ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama error status %d: %s", resp.StatusCode, string(data))
	}

	var oResp ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&oResp); err != nil {
		return nil, fmt.Errorf("decode ollama response: %w", err)
	}

	// Clean the response by removing markdown code blocks if present
	response := oResp.Response
	// Remove ```json and ``` markers if they exist
	if len(response) >= 7 && response[:7] == "```json" {
		response = response[7:]
	}
	if len(response) >= 3 && response[len(response)-3:] == "```" {
		response = response[:len(response)-3]
	}

	// The model's textual response should itself be JSON describing a Ticket.
	var ticket models.Ticket
	if err := json.Unmarshal([]byte(response), &ticket); err != nil {
		// Log the problematic response for debugging
		log.Printf("[DEBUG] Failed to unmarshal LLM response. Response was: %s", response)
		return nil, fmt.Errorf("unmarshal ticket JSON from llm: %w (response: %s)", err, response)
	}

	return &ticket, nil
}
