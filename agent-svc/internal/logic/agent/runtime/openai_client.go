package runtime

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type openAICompatibleModelClient struct {
	baseURL string
	apiKey  string
	model   string
	client  *http.Client
}

type openAIChatCompletionRequest struct {
	Model          string                    `json:"model"`
	Messages       []openAIChatMessage       `json:"messages"`
	Temperature    float64                   `json:"temperature,omitempty"`
	ResponseFormat *openAIChatResponseFormat `json:"response_format,omitempty"`
}

type openAIChatResponseFormat struct {
	Type string `json:"type"`
}

type openAIChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIChatCompletionResponse struct {
	Choices []struct {
		Message openAIChatMessage `json:"message"`
	} `json:"choices"`
}

func NewOpenAICompatibleModelClient(baseURL, apiKey, model string, client *http.Client) ModelClient {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	apiKey = strings.TrimSpace(apiKey)
	model = strings.TrimSpace(model)
	if baseURL == "" || apiKey == "" || model == "" {
		return nil
	}
	if client == nil {
		client = &http.Client{Timeout: 20 * time.Second}
	}
	return &openAICompatibleModelClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		model:   model,
		client:  client,
	}
}

func (c *openAICompatibleModelClient) Complete(ctx context.Context, req ModelCompletionRequest) (string, error) {
	payload, err := json.Marshal(openAIChatCompletionRequest{
		Model: c.model,
		Messages: []openAIChatMessage{
			{Role: "system", Content: strings.TrimSpace(req.SystemPrompt)},
			{Role: "user", Content: strings.TrimSpace(req.UserPrompt)},
		},
		Temperature: 0.1,
		ResponseFormat: &openAIChatResponseFormat{
			Type: "json_object",
		},
	})
	if err != nil {
		return "", err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return "", &OrderLookupError{Code: orderLookupErrorDownstream}
	}

	var decoded openAIChatCompletionResponse
	if err = json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return "", err
	}
	if len(decoded.Choices) == 0 {
		return "", nil
	}
	return strings.TrimSpace(decoded.Choices[0].Message.Content), nil
}
