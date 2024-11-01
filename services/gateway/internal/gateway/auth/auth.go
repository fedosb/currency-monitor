package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fedosb/currency-monitor/services/gateway/internal/config"
)

type Gateway struct {
	baseUrl string
	client  *http.Client
}

func New(cfg config.AuthApiConfig) *Gateway {
	return &Gateway{baseUrl: cfg.GetURL(), client: &http.Client{}}
}

func (g *Gateway) GenerateToken(ctx context.Context, login string) (string, error) {
	path := fmt.Sprintf("/generate?login=%s", login)
	token, err := g.request(ctx, path, nil)
	if err != nil {
		return "", fmt.Errorf("request auth api: %w", err)
	}

	return token, nil
}

func (g *Gateway) ValidateToken(ctx context.Context, token string) error {
	path := "/validate"
	_, err := g.request(ctx, path, map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)})
	if err != nil {
		return fmt.Errorf("request auth api: %w", err)
	}

	return nil
}

func (g *Gateway) request(ctx context.Context, path string, headers map[string]string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, g.baseUrl+path, nil)
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	var response string
	if resp.ContentLength != 0 {
		if _, err := fmt.Fscan(resp.Body, &response); err != nil {
			return "", fmt.Errorf("scan response: %w", err)
		}
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("got %d with response: %s", resp.StatusCode, response)
	}

	return response, nil
}
