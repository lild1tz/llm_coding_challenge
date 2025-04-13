package apollo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/models"
)

func newClient(cfg Config) *client {
	return &client{
		Client: &http.Client{},
		cfg:    cfg,
	}
}

type client struct {
	*http.Client
	cfg Config
}

func (c *client) Release() error {
	return nil
}

type RequestBody struct {
	Message string `json:"message"`
}

type ResponseBody struct {
	Table models.Table `json:"table"`
}

func (c *client) PredictTableFromText(ctx context.Context, text string) (models.Table, error) {
	jsonBody, err := json.Marshal(RequestBody{Message: text})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", c.cfg.ApolloURL+"/process_message", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}

	var responseBody ResponseBody
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	defer resp.Body.Close()

	return responseBody.Table, nil
}
