package openai

import (
	"bytes"
	"cold_emailer/constants"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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

// EmailGenerationData contains all the data needed for personalized email generation
type EmailGenerationData struct {
	ContactName       string
	ContactRole       string
	ContactLinkedIn   string
	CompanyName       string
	CompanyWebsite    string
	CompanyIndustry   string
	ProfileName       string
	ProfileExperience string
	ProfileSkills     string
	ProfileSummary    string
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

// GeneratePersonalizedEmail generates a personalized email for a specific contact
func (c *OpenAIClient) GeneratePersonalizedEmail(data EmailGenerationData) (subject, body string, err error) {
	prompt := fmt.Sprintf(`
Generate a personalized cold outreach email for a job opportunity. 

CONTACT INFORMATION:
- Name: %s
- Role: %s
- LinkedIn: %s

COMPANY INFORMATION:
- Company Name: %s
- Website: %s
- Industry: %s

MY PROFILE:
- Name: %s
- Experience: %s
- Skills: %s
- Summary: %s

INSTRUCTIONS:
1. Write a professional, personalized email that shows I've researched their company
2. Mention their specific role and company
3. Highlight relevant skills and experience that match their industry
4. Keep the email concise but impactful (150-200 words)
5. Include a clear call-to-action
6. Make it sound natural and not overly salesy
7. Show enthusiasm for their company and role

Generate both a compelling subject line and the email body.
`,
		data.ContactName, data.ContactRole, data.ContactLinkedIn,
		data.CompanyName, data.CompanyWebsite, data.CompanyIndustry,
		data.ProfileName, data.ProfileExperience, data.ProfileSkills, data.ProfileSummary)

	url := "https://api.openai.com/v1/chat/completions"

	messages := []ChatMessage{
		{Role: "system", Content: "You are a professional email writer. Generate a subject line and email body. Format your response exactly as: SUBJECT: [subject line] ---BODY--- [email body only, no headers or labels]"},
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
		return "", "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to call OpenAI API: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("OpenAI API error: %s", string(respBody))
	}

	var completionResp ChatCompletionResponse
	if err := json.Unmarshal(respBody, &completionResp); err != nil {
		return "", "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if completionResp.Error != nil {
		return "", "", fmt.Errorf("OpenAI error: %s", completionResp.Error.Message)
	}

	if len(completionResp.Choices) == 0 {
		return "", "", fmt.Errorf("no choices returned from OpenAI")
	}

	content := completionResp.Choices[0].Message.Content

	// Parse subject and body from the response
	parts := strings.Split(content, "---BODY---")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid response format from OpenAI: expected '---BODY---' separator")
	}

	subjectPart := strings.TrimSpace(parts[0])
	bodyPart := strings.TrimSpace(parts[1])

	// Extract subject (remove "SUBJECT:" prefix)
	if strings.HasPrefix(subjectPart, "SUBJECT:") {
		subject = strings.TrimSpace(strings.Replace(subjectPart, "SUBJECT:", "", 1))
	} else {
		subject = subjectPart
	}

	// Clean up body (remove any remaining headers or labels)
	body = bodyPart
	// Remove common prefixes that might be included
	body = strings.TrimPrefix(body, "Email body:")
	body = strings.TrimPrefix(body, "Body:")
	body = strings.TrimPrefix(body, "EMAIL:")
	body = strings.TrimSpace(body)

	return subject, body, nil
}
