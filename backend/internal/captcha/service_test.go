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
	svc := NewService(fakeVerifyClient{result: VerifyResult{Success: true, Score: 0.9}}, 0.5, true)
	if err := svc.Verify(context.Background(), "token"); err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
}

func TestServiceVerifyRejectLowScore(t *testing.T) {
	svc := NewService(fakeVerifyClient{result: VerifyResult{Success: true, Score: 0.4}}, 0.5, true)
	err := svc.Verify(context.Background(), "token")
	if !errors.Is(err, ErrCaptchaFailed) {
		t.Fatalf("Verify() error = %v, want %v", err, ErrCaptchaFailed)
	}
}

func TestServiceVerifyDisabledBypassesValidation(t *testing.T) {
	svc := NewService(fakeVerifyClient{err: errors.New("should not call")}, 0.5, false)
	if err := svc.Verify(context.Background(), ""); err != nil {
		t.Fatalf("Verify() error = %v, want nil", err)
	}
}
