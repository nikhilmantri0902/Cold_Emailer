package openai

import (
	"bytes"
	"cold_emailer/constants"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// OpenAIClient handles requests to the OpenAI API
type OpenAIClient struct {
	APIKey              string
	Model               string
	Temperature         float32
	MaxCompletionTokens int
}

// NewOpenAIClient creates a new OpenAI client with config from constants
func NewOpenAIClient() *OpenAIClient {
	return &OpenAIClient{
		APIKey:              constants.OPENAI_API_KEY,
		Model:               constants.OPENAI_MODEL,
		Temperature:         constants.OPENAI_TEMPERATURE,
		MaxCompletionTokens: constants.OPENAI_MAX_COMPLETION_TOKENS,
	}
}

// ChatMessage represents a message for the OpenAI chat API
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionRequest is the request body for OpenAI chat API
type ChatCompletionRequest struct {
	Model               string        `json:"model"`
	Messages            []ChatMessage `json:"messages"`
	Temperature         float32       `json:"temperature"`
	MaxCompletionTokens int           `json:"max_completion_tokens"`
}

// ChatCompletionResponse is the response from OpenAI chat API
type ChatCompletionResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// GenerateEmail calls OpenAI API with the given prompt and returns the generated email text
func (c *OpenAIClient) GenerateEmail(prompt string) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"

	messages := []ChatMessage{
		{Role: "system", Content: constants.PROMPT_EMAIL_GENERATION},
		{Role: "user", Content: prompt},
	}

	reqBody := ChatCompletionRequest{
		Model:               c.Model,
		Messages:            messages,
		Temperature:         c.Temperature,
		MaxCompletionTokens: c.MaxCompletionTokens,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Println("error:", err)
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Println("error:", err)
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error:", err)
		return "", fmt.Errorf("failed to call OpenAI API: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error:", err)
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != 200 {
		log.Println("error: OpenAI API error", string(respBody))
		return "", fmt.Errorf("OpenAI API error: %s", string(respBody))
	}

	var completionResp ChatCompletionResponse
	if err := json.Unmarshal(respBody, &completionResp); err != nil {
		log.Println("error:", err)
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if completionResp.Error != nil {
		log.Println("error: OpenAI error", completionResp.Error.Message)
		return "", fmt.Errorf("OpenAI error: %s", completionResp.Error.Message)
	}

	if len(completionResp.Choices) == 0 {
		log.Println("error: no choices returned from OpenAI")
		return "", fmt.Errorf("no choices returned from OpenAIII")
	}

	return completionResp.Choices[0].Message.Content, nil
}
