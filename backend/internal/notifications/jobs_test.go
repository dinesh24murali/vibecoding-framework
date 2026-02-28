package notifications

import (
	"testing"
	"time"
)

func TestIsRetryDue(t *testing.T) {
	now := time.Now().UTC()

	if !IsRetryDue("queued", 0, now, now) {
		t.Fatalf("queued job should be due")
	}

	updatedAt := now.Add(-11 * time.Second)
	if !IsRetryDue("failed", 1, updatedAt, now) {
		t.Fatalf("failed job with first backoff elapsed should be due")
	}

	updatedAt = now.Add(-5 * time.Second)
	if IsRetryDue("failed", 1, updatedAt, now) {
		t.Fatalf("failed job before backoff should not be due")
	}
}
