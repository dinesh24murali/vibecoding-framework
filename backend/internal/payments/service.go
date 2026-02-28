package payments

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dinesh/vibecoding-framework/backend/internal/captcha"
	"gorm.io/gorm"
)

var (
	ErrPaymentNotFound  = errors.New("payment not found")
	ErrServiceNotFound  = errors.New("service request not found")
	ErrCaptchaRejected  = errors.New("Please try after some time")
	ErrInvalidSignature = errors.New("invalid gateway signature")
)

type captchaVerifier interface {
	Verify(ctx context.Context, token string) error
}

type Service struct {
	db         *gorm.DB
	client     *RazorpayClient
	retryAudit *RetryAuditRepo
	captcha    captchaVerifier
}

type RetryResult struct {
	ServiceRequest map[string]any `json:"service_request"`
	Payment        map[string]any `json:"payment"`
}

func NewService(db *gorm.DB, client *RazorpayClient, retryAudit *RetryAuditRepo, captchaSvc captchaVerifier) *Service {
	return &Service{db: db, client: client, retryAudit: retryAudit, captcha: captchaSvc}
}

func (s *Service) RetryPayment(ctx context.Context, serviceRequestID int64, recaptchaToken string, idempotencyKey string) (*RetryResult, error) {
	if err := s.captcha.Verify(ctx, recaptchaToken); err != nil {
		if errors.Is(err, captcha.ErrCaptchaFailed) {
			return nil, ErrCaptchaRejected
		}
		return nil, err
	}

	var req struct {
		ID       int64   `gorm:"column:id"`
		PlanID   int64   `gorm:"column:plan_id"`
		Status   string  `gorm:"column:status"`
		Price    float64 `gorm:"column:price"`
		Discount float64 `gorm:"column:discount"`
	}
	err := s.db.WithContext(ctx).Table("service_requests").
		Select("service_requests.id, service_requests.plan_id, service_requests.status, plans.price, plans.discount").
		Joins("JOIN plans ON plans.id = service_requests.plan_id").
		Where("service_requests.id = ? AND service_requests.deleted_at IS NULL", serviceRequestID).
		Take(&req).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrServiceNotFound
		}
		return nil, err
	}

	amount := req.Price - req.Discount
	if amount < 0.01 {
		amount = 0.01
	}

	orderID, err := s.client.CreateOrder(ctx, amount, "INR", fmt.Sprintf("sr_%d", serviceRequestID))
	if err != nil {
		return nil, err
	}

	payload := map[string]any{
		"service_request_id": serviceRequestID,
		"gateway_order_id":   orderID,
		"amount":             amount,
		"status":             "created",
		"idempotency_key":    idempotencyKey,
	}
	if err := s.db.WithContext(ctx).Table("payment_attempts").Create(payload).Error; err != nil {
		return nil, fmt.Errorf("create payment attempt: %w", err)
	}

	retryNo, err := s.retryAudit.NextRetryNo(ctx, serviceRequestID)
	if err != nil {
		return nil, err
	}
	if err := s.retryAudit.Insert(ctx, serviceRequestID, retryNo, "customer_retry"); err != nil {
		return nil, err
	}

	return &RetryResult{
		ServiceRequest: map[string]any{
			"id":     fmt.Sprintf("%d", serviceRequestID),
			"status": req.Status,
		},
		Payment: map[string]any{
			"gateway":  "razorpay",
			"order_id": orderID,
			"amount":   amount,
			"currency": "INR",
			"status":   "created",
		},
	}, nil
}

func (s *Service) GetCustomerPaymentStatus(ctx context.Context, serviceRequestID int64, phone string) (map[string]any, error) {
	var row struct {
		ID        int64     `gorm:"column:id"`
		UserID    int64     `gorm:"column:user_id"`
		PlanID    int64     `gorm:"column:plan_id"`
		Status    string    `gorm:"column:status"`
		CreatedAt time.Time `gorm:"column:created_at"`
		UpdatedAt time.Time `gorm:"column:updated_at"`
	}
	err := s.db.WithContext(ctx).Table("service_requests").
		Select("service_requests.id, service_requests.user_id, service_requests.plan_id, service_requests.status, service_requests.created_at, service_requests.updated_at").
		Joins("JOIN users ON users.id = service_requests.user_id").
		Where("service_requests.id = ? AND users.phone_number = ?", serviceRequestID, phone).
		Take(&row).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrServiceNotFound
		}
		return nil, err
	}

	return map[string]any{
		"service_request": map[string]any{
			"id":         fmt.Sprintf("%d", row.ID),
			"user_id":    fmt.Sprintf("%d", row.UserID),
			"plan_id":    fmt.Sprintf("%d", row.PlanID),
			"status":     row.Status,
			"created_at": row.CreatedAt.UTC().Format(time.RFC3339),
			"updated_at": row.UpdatedAt.UTC().Format(time.RFC3339),
		},
	}, nil
}

func (s *Service) HandleCallback(ctx context.Context, rawBody []byte, signature string, body map[string]any) (bool, error) {
	if !s.client.VerifyWebhookSignature(rawBody, signature) {
		return false, ErrInvalidSignature
	}

	orderID, paymentID, status := extractPaymentFields(body)
	if orderID == "" {
		return false, ErrPaymentNotFound
	}

	var attempt struct {
		ID               int64   `gorm:"column:id"`
		ServiceRequestID int64   `gorm:"column:service_request_id"`
		GatewayPaymentID *string `gorm:"column:gateway_payment_id"`
	}
	err := s.db.WithContext(ctx).Table("payment_attempts").Where("gateway_order_id = ?", orderID).Take(&attempt).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, ErrPaymentNotFound
		}
		return false, err
	}

	if attempt.GatewayPaymentID != nil && paymentID != "" && *attempt.GatewayPaymentID == paymentID {
		return true, nil
	}

	callbackJSON, _ := json.Marshal(body)
	updates := map[string]any{"status": status, "callback_payload": string(callbackJSON)}
	if paymentID != "" {
		updates["gateway_payment_id"] = paymentID
	}
	if err := s.db.WithContext(ctx).Table("payment_attempts").Where("id = ?", attempt.ID).Updates(updates).Error; err != nil {
		return false, err
	}

	serviceStatus := "payment_failed"
	if status == "captured" || status == "authorized" {
		serviceStatus = "payment_success"
	}
	if err := s.db.WithContext(ctx).Table("service_requests").Where("id = ?", attempt.ServiceRequestID).Update("status", serviceStatus).Error; err != nil {
		return false, err
	}

	return false, nil
}

func extractPaymentFields(body map[string]any) (orderID string, paymentID string, status string) {
	status = "failed"
	payload, _ := body["payload"].(map[string]any)
	payment, _ := payload["payment"].(map[string]any)
	entity, _ := payment["entity"].(map[string]any)

	if v, ok := entity["order_id"].(string); ok {
		orderID = v
	}
	if v, ok := entity["id"].(string); ok {
		paymentID = v
	}
	if v, ok := entity["status"].(string); ok && v != "" {
		status = v
	}

	return orderID, paymentID, status
}
