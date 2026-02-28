package captcha

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var ErrCaptchaFailed = errors.New("Please try after some time")

type Service struct {
	client   VerifyClient
	minScore float64
}

func NewService(client VerifyClient, minScore float64) *Service {
	return &Service{client: client, minScore: minScore}
}

func (s *Service) Verify(ctx context.Context, token string) error {
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
