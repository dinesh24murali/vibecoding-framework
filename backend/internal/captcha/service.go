package captcha

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var ErrCaptchaFailed = errors.New("Please try after some time")

type Service struct {
	enabled  bool
	client   VerifyClient
	minScore float64
}

func NewService(client VerifyClient, minScore float64, enabled bool) *Service {
	return &Service{enabled: enabled, client: client, minScore: minScore}
}

func (s *Service) Verify(ctx context.Context, token string) error {
	if !s.enabled {
		return nil
	}

	if strings.TrimSpace(token) == "" {
		return ErrCaptchaFailed
	}

	result, err := s.client.Verify(ctx, token)
	if err != nil {
		return fmt.Errorf("verify captcha: %w", err)
	}

	if !result.Success || result.Score < s.minScore {
		return ErrCaptchaFailed
	}

	return nil
}
