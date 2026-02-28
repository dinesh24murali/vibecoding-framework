package notifications

import "time"

const CompletedTemplate = "service request %d is completed"

func RetryBackoff(failureCount int) time.Duration {
	switch failureCount {
	case 1:
		return 10 * time.Second
	case 2:
		return 30 * time.Second
	case 3:
		return 2 * time.Minute
	default:
		return 0
	}
}

func IsRetryDue(status string, failureCount int, updatedAt time.Time, now time.Time) bool {
	if status == "queued" {
		return true
	}

	if status != "failed" {
		return false
	}

	backoff := RetryBackoff(failureCount)
	if backoff == 0 {
		return false
	}

	return !now.Before(updatedAt.Add(backoff))
}
