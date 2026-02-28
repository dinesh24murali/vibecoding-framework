package servicerequests

import "errors"

const (
	StatusPending        = "pending"
	StatusPaymentFailed  = "payment_failed"
	StatusPaymentSuccess = "payment_success"
	StatusBlocked        = "blocked"
	StatusCompleted      = "completed"
)

var ErrInvalidStatus = errors.New("invalid service request status")

var allowedTransitions = map[string]map[string]struct{}{
	StatusPending: {
		StatusPaymentFailed:  {},
		StatusPaymentSuccess: {},
		StatusBlocked:        {},
	},
	StatusPaymentFailed: {
		StatusPaymentSuccess: {},
		StatusBlocked:        {},
	},
	StatusPaymentSuccess: {
		StatusBlocked:   {},
		StatusCompleted: {},
	},
	StatusBlocked: {
		StatusCompleted: {},
	},
	StatusCompleted: {},
}

func IsValidStatus(status string) bool {
	_, ok := allowedTransitions[status]
	return ok
}

func CanTransition(from string, to string) (bool, error) {
	if !IsValidStatus(from) || !IsValidStatus(to) {
		return false, ErrInvalidStatus
	}

	if from == to {
		return true, nil
	}

	_, ok := allowedTransitions[from][to]
	return ok, nil
}
