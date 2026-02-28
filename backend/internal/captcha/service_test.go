package captcha

import (
	"context"
	"errors"
	"testing"
)

type fakeVerifyClient struct {
	result VerifyResult
	err    error
}

func (f fakeVerifyClient) Verify(_ context.Context, _ string) (VerifyResult, error) {
	if f.err != nil {
		return VerifyResult{}, f.err
	}

	return f.result, nil
}

func TestServiceVerifyPass(t *testing.T) {
	svc := NewService(fakeVerifyClient{result: VerifyResult{Success: true, Score: 0.9}}, 0.5)
	if err := svc.Verify(context.Background(), "token"); err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
}

func TestServiceVerifyRejectLowScore(t *testing.T) {
	svc := NewService(fakeVerifyClient{result: VerifyResult{Success: true, Score: 0.4}}, 0.5)
	err := svc.Verify(context.Background(), "token")
	if !errors.Is(err, ErrCaptchaFailed) {
		t.Fatalf("Verify() error = %v, want %v", err, ErrCaptchaFailed)
	}
}
