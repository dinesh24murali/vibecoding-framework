package servicerequests

import "testing"

func TestCanTransition(t *testing.T) {
	ok, err := CanTransition(StatusPaymentSuccess, StatusCompleted)
	if err != nil {
		t.Fatalf("CanTransition() error = %v", err)
	}
	if !ok {
		t.Fatalf("CanTransition() = false, want true")
	}

	ok, err = CanTransition(StatusPending, StatusCompleted)
	if err != nil {
		t.Fatalf("CanTransition() error = %v", err)
	}
	if ok {
		t.Fatalf("CanTransition() = true, want false")
	}
}
