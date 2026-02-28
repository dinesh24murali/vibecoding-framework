package checkout

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dinesh/vibecoding-framework/backend/internal/captcha"
	"github.com/dinesh/vibecoding-framework/backend/internal/plans"
	"github.com/dinesh/vibecoding-framework/backend/internal/servicerequests"
	"gorm.io/gorm"
)

var (
	ErrInvalidInput  = errors.New("invalid checkout input")
	ErrNotFound      = errors.New("resource not found")
	ErrConflict      = errors.New("checkout conflict")
	ErrCaptchaFailed = errors.New("Please try after some time")
)

type captchaVerifier interface {
	Verify(ctx context.Context, token string) error
}

type paymentOrderCreator interface {
	CreateOrder(ctx context.Context, amount float64, currency string, receipt string) (string, error)
}

type Service struct {
	db                 *gorm.DB
	plansRepo          *plans.Repository
	serviceRequestRepo *servicerequests.Repository
	captcha            captchaVerifier
	payments           paymentOrderCreator
}

type CreateCheckoutInput struct {
	ProviderID     int64
	PlanID         int64
	CustomerName   string
	CustomerPhone  string
	RecaptchaToken string
	IdempotencyKey string
}

type CreateCheckoutResult struct {
	ServiceRequest ServiceRequestPayload `json:"service_request"`
	Payment        PaymentPayload        `json:"payment"`
}

type ServiceRequestPayload struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	PlanID    string `json:"plan_id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type PaymentPayload struct {
	Gateway  string  `json:"gateway"`
	OrderID  string  `json:"order_id"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Status   string  `json:"status"`
}

type userRow struct {
	ID          int64  `gorm:"column:id"`
	Name        string `gorm:"column:name"`
	PhoneNumber string `gorm:"column:phone_number"`
	Role        string `gorm:"column:role"`
}

func (userRow) TableName() string {
	return "users"
}

func NewService(db *gorm.DB, plansRepo *plans.Repository, serviceRequestRepo *servicerequests.Repository, captchaService captchaVerifier, paymentService paymentOrderCreator) *Service {
	return &Service{db: db, plansRepo: plansRepo, serviceRequestRepo: serviceRequestRepo, captcha: captchaService, payments: paymentService}
}

func (s *Service) CreateServiceRequestAndPayment(ctx context.Context, input CreateCheckoutInput) (*CreateCheckoutResult, error) {
	if err := validateInput(input); err != nil {
		return nil, err
	}

	if err := s.captcha.Verify(ctx, input.RecaptchaToken); err != nil {
		if errors.Is(err, captcha.ErrCaptchaFailed) {
			return nil, ErrCaptchaFailed
		}

		return nil, fmt.Errorf("verify recaptcha: %w", err)
	}

	plan, err := s.plansRepo.FindActiveByID(ctx, input.PlanID)
	if err != nil {
		return nil, err
	}
	if plan == nil || plan.ProviderID != input.ProviderID {
		return nil, ErrNotFound
	}

	userID, err := s.upsertCustomer(ctx, input.CustomerName, input.CustomerPhone)
	if err != nil {
		return nil, err
	}

	serviceRequest, err := s.serviceRequestRepo.CreatePending(ctx, userID, plan.ID)
	if err != nil {
		return nil, err
	}

	amount := plan.Price - plan.Discount
	if amount < 0.01 {
		amount = 0.01
	}

	orderID, err := s.payments.CreateOrder(ctx, amount, "INR", fmt.Sprintf("sr_%d", serviceRequest.ID))
	if err != nil {
		return nil, fmt.Errorf("create razorpay order: %w", err)
	}
	if err := s.createPaymentAttempt(ctx, serviceRequest.ID, orderID, amount, input.IdempotencyKey); err != nil {
		return nil, err
	}

	return &CreateCheckoutResult{
		ServiceRequest: ServiceRequestPayload{
			ID:        fmt.Sprintf("%d", serviceRequest.ID),
			UserID:    fmt.Sprintf("%d", serviceRequest.UserID),
			PlanID:    fmt.Sprintf("%d", serviceRequest.PlanID),
			Status:    serviceRequest.Status,
			CreatedAt: serviceRequest.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt: serviceRequest.UpdatedAt.UTC().Format(time.RFC3339),
		},
		Payment: PaymentPayload{
			Gateway:  "razorpay",
			OrderID:  orderID,
			Amount:   amount,
			Currency: "INR",
			Status:   "created",
		},
	}, nil
}

func validateInput(input CreateCheckoutInput) error {
	if input.ProviderID <= 0 || input.PlanID <= 0 {
		return ErrInvalidInput
	}
	if strings.TrimSpace(input.CustomerName) == "" || strings.TrimSpace(input.CustomerPhone) == "" {
		return ErrInvalidInput
	}
	if strings.TrimSpace(input.RecaptchaToken) == "" {
		return ErrInvalidInput
	}
	if strings.TrimSpace(input.IdempotencyKey) == "" {
		return ErrInvalidInput
	}

	return nil
}

func (s *Service) upsertCustomer(ctx context.Context, name string, phone string) (int64, error) {
	cleanName := strings.TrimSpace(name)
	cleanPhone := strings.TrimSpace(phone)

	var existing userRow
	err := s.db.WithContext(ctx).Where("phone_number = ?", cleanPhone).First(&existing).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, fmt.Errorf("query customer by phone: %w", err)
	}

	if err == nil {
		if existing.Role == "admin" {
			return 0, ErrConflict
		}

		if existing.Name != cleanName {
			if updateErr := s.db.WithContext(ctx).Model(&userRow{}).Where("id = ?", existing.ID).Update("name", cleanName).Error; updateErr != nil {
				return 0, fmt.Errorf("update customer name: %w", updateErr)
			}
		}

		return existing.ID, nil
	}

	customer := userRow{Name: cleanName, PhoneNumber: cleanPhone, Role: "customer"}
	if createErr := s.db.WithContext(ctx).Table("users").Create(&customer).Error; createErr != nil {
		return 0, fmt.Errorf("create customer user: %w", createErr)
	}

	return customer.ID, nil
}

func (s *Service) createPaymentAttempt(ctx context.Context, serviceRequestID int64, orderID string, amount float64, idempotencyKey string) error {
	payload := map[string]any{
		"service_request_id": serviceRequestID,
		"gateway_order_id":   orderID,
		"amount":             amount,
		"status":             "created",
		"idempotency_key":    idempotencyKey,
	}

	if err := s.db.WithContext(ctx).Table("payment_attempts").Create(payload).Error; err != nil {
		return fmt.Errorf("create payment attempt: %w", err)
	}

	return nil
}
