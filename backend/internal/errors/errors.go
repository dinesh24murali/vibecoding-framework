package errors

type ErrorDetail struct {
	Field string `json:"field"`
	Issue string `json:"issue"`
}

type ErrorBody struct {
	Code      string        `json:"code"`
	Message   string        `json:"message"`
	Details   []ErrorDetail `json:"details,omitempty"`
	RequestID string        `json:"request_id,omitempty"`
}

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

func NewInternal(requestID string) ErrorResponse {
	return ErrorResponse{
		Error: ErrorBody{
			Code:      "INTERNAL_ERROR",
			Message:   "internal server error",
			RequestID: requestID,
		},
	}
}
