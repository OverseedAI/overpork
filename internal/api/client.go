package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/OverseedAI/overpork/internal/config"
)

const BaseURL = "https://api.porkbun.com/api/json/v3"

type Client struct {
	httpClient *http.Client
	apiKey     string
	secretKey  string
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		apiKey:     cfg.APIKey,
		secretKey:  cfg.SecretKey,
	}
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func (c *Client) doURL(method, url string, reqBody, respBody any) error {
	var body io.Reader
	if reqBody != nil {
		data, err := json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}
		body = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if respBody != nil {
		if err := json.Unmarshal(respData, respBody); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
	}

	// Check for API error
	var baseResp Response
	if err := json.Unmarshal(respData, &baseResp); err == nil {
		if baseResp.Status == "ERROR" {
			return fmt.Errorf("API error: %s", baseResp.Message)
		}
	}

	return nil
}

func (c *Client) post(endpoint string, reqBody, respBody any) error {
	return c.doURL("POST", BaseURL+endpoint, reqBody, respBody)
}

func (c *Client) postURL(url string, reqBody, respBody any) error {
	return c.doURL("POST", url, reqBody, respBody)
}

func (c *Client) authBody() map[string]string {
	return map[string]string{
		"apikey":       c.apiKey,
		"secretapikey": c.secretKey,
	}
}

func (c *Client) authBodyWith(extra map[string]any) map[string]any {
	body := map[string]any{
		"apikey":       c.apiKey,
		"secretapikey": c.secretKey,
	}
	for k, v := range extra {
		body[k] = v
	}
	return body
}

// Ping checks API connectivity
func (c *Client) Ping() error {
	var resp Response
	return c.post("/ping", c.authBody(), &resp)
}
