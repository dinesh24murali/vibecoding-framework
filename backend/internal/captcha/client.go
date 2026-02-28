package captcha

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const googleVerifyURL = "https://www.google.com/recaptcha/api/siteverify"

type VerifyClient interface {
	Verify(ctx context.Context, token string) (VerifyResult, error)
}

type GoogleClient struct {
	httpClient *http.Client
	secretKey  string
}

type VerifyResult struct {
	Success bool
	Score   float64
}

type googleVerifyResponse struct {
	Success bool    `json:"success"`
	Score   float64 `json:"score"`
}

func NewGoogleClient(secretKey string) *GoogleClient {
	return &GoogleClient{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		secretKey:  secretKey,
	}
}

func (c *GoogleClient) Verify(ctx context.Context, token string) (VerifyResult, error) {
	values := url.Values{}
	values.Set("secret", c.secretKey)
	values.Set("response", strings.TrimSpace(token))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, googleVerifyURL, strings.NewReader(values.Encode()))
	if err != nil {
		return VerifyResult{}, fmt.Errorf("create captcha verify request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return VerifyResult{}, fmt.Errorf("execute captcha verify request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return VerifyResult{}, fmt.Errorf("captcha verify returned status %d", resp.StatusCode)
	}

	var payload googleVerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return VerifyResult{}, fmt.Errorf("decode captcha verify response: %w", err)
	}

	return VerifyResult{Success: payload.Success, Score: payload.Score}, nil
}
