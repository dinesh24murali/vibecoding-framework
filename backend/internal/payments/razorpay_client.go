package payments

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type RazorpayClient struct {
	keyID         string
	keySecret     string
	webhookSecret string
}

func NewRazorpayClient(keyID string, keySecret string, webhookSecret string) *RazorpayClient {
	return &RazorpayClient{keyID: keyID, keySecret: keySecret, webhookSecret: webhookSecret}
}

func (c *RazorpayClient) CreateOrder(_ context.Context, amount float64, _ string, receipt string) (string, error) {
	if amount <= 0 {
		return "", fmt.Errorf("invalid amount")
	}

	if receipt == "" {
		receipt = "dth"
	}

	return fmt.Sprintf("order_%s_%d", receipt, time.Now().UnixNano()), nil
}

func (c *RazorpayClient) VerifyWebhookSignature(rawBody []byte, signature string) bool {
	if signature == "" || c.webhookSecret == "" {
		return false
	}

	h := hmac.New(sha256.New, []byte(c.webhookSecret))
	_, _ = h.Write(rawBody)
	computed := hex.EncodeToString(h.Sum(nil))
	return hmac.Equal([]byte(computed), []byte(signature))
}
