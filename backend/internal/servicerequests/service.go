package servicerequests

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

var (
	ErrServiceRequestNotFound = errors.New("service request not found")
	ErrServiceRequestInvalid  = errors.New("invalid service request input")
	ErrInvalidTransition      = errors.New("invalid service request transition")
)

type Service struct {
	db *gorm.DB
}

type ListAdminInput struct {
	Page       int
	PageSize   int
	Status     string
	ProviderID *int64
	FromDate   string
	ToDate     string
	Sort       string
}

type ServiceRequestResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	PlanID    string `json:"plan_id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserSummary struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Phone     string `json:"phone_number"`
	Role      string `json:"role"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ProviderSummary struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	ImageURL  *string `json:"image_url,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type PlanSummary struct {
	ID          string  `json:"id"`
	ProviderID  string  `json:"provider_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Discount    float64 `json:"discount"`
	IsActive    bool    `json:"is_active"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type ServiceRequestDetailsResponse struct {
	ServiceRequestResponse
	User     UserSummary     `json:"user"`
	Plan     PlanSummary     `json:"plan"`
	Provider ProviderSummary `json:"provider"`
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) ListAdmin(ctx context.Context, input ListAdminInput) ([]ServiceRequestResponse, int64, error) {
	page, pageSize, err := normalizePaging(input.Page, input.PageSize)
	if err != nil {
		return nil, 0, err
	}

	query := s.db.WithContext(ctx).
		Table("service_requests").
		Joins("JOIN plans ON plans.id = service_requests.plan_id").
		Where("service_requests.deleted_at IS NULL")

	status := strings.TrimSpace(input.Status)
	if status != "" {
		if !IsValidStatus(status) {
			return nil, 0, ErrServiceRequestInvalid
		}
		query = query.Where("service_requests.status = ?", status)
	}

	if input.ProviderID != nil {
		if *input.ProviderID <= 0 {
			return nil, 0, ErrServiceRequestInvalid
		}
		query = query.Where("plans.provider_id = ?", *input.ProviderID)
	}

	if input.FromDate != "" {
		fromDate, parseErr := time.Parse("2006-01-02", input.FromDate)
		if parseErr != nil {
			return nil, 0, ErrServiceRequestInvalid
		}
		query = query.Where("service_requests.created_at >= ?", fromDate)
	}

	if input.ToDate != "" {
		toDate, parseErr := time.Parse("2006-01-02", input.ToDate)
		if parseErr != nil {
			return nil, 0, ErrServiceRequestInvalid
		}
		query = query.Where("service_requests.created_at < ?", toDate.Add(24*time.Hour))
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count service requests: %w", err)
	}

	var rows []ServiceRequest
	if err := query.
		Order(resolveSort(input.Sort)).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("list service requests: %w", err)
	}

	items := make([]ServiceRequestResponse, 0, len(rows))
	for _, row := range rows {
		items = append(items, toServiceRequestResponse(row))
	}

	return items, total, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*ServiceRequestDetailsResponse, error) {
	if id <= 0 {
		return nil, ErrServiceRequestInvalid
	}

	serviceRequest, err := s.repoGetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	var user struct {
		ID        int64     `gorm:"column:id"`
		Name      string    `gorm:"column:name"`
		Phone     string    `gorm:"column:phone_number"`
		Role      string    `gorm:"column:role"`
		Status    string    `gorm:"column:status"`
		CreatedAt time.Time `gorm:"column:created_at"`
		UpdatedAt time.Time `gorm:"column:updated_at"`
	}
	if err := s.db.WithContext(ctx).
		Table("users").
		Select("id, name, phone_number, role, status, created_at, updated_at").
		Where("id = ?", serviceRequest.UserID).
		Take(&user).Error; err != nil {
		return nil, fmt.Errorf("get service request user: %w", err)
	}

	var plan struct {
		ID          int64     `gorm:"column:id"`
		ProviderID  int64     `gorm:"column:provider_id"`
		Name        string    `gorm:"column:name"`
		Description string    `gorm:"column:description"`
		Price       float64   `gorm:"column:price"`
		Discount    float64   `gorm:"column:discount"`
		IsActive    bool      `gorm:"column:is_active"`
		CreatedAt   time.Time `gorm:"column:created_at"`
		UpdatedAt   time.Time `gorm:"column:updated_at"`
	}
	if err := s.db.WithContext(ctx).
		Table("plans").
		Select("id, provider_id, name, description, price, discount, is_active, created_at, updated_at").
		Where("id = ?", serviceRequest.PlanID).
		Take(&plan).Error; err != nil {
		return nil, fmt.Errorf("get service request plan: %w", err)
	}

	var provider struct {
		ID        int64     `gorm:"column:id"`
		Name      string    `gorm:"column:name"`
		ImageURL  *string   `gorm:"column:image_url"`
		CreatedAt time.Time `gorm:"column:created_at"`
		UpdatedAt time.Time `gorm:"column:updated_at"`
	}
	if err := s.db.WithContext(ctx).
		Table("providers").
		Select("id, name, image_url, created_at, updated_at").
		Where("id = ?", plan.ProviderID).
		Take(&provider).Error; err != nil {
		return nil, fmt.Errorf("get service request provider: %w", err)
	}

	response := &ServiceRequestDetailsResponse{
		ServiceRequestResponse: toServiceRequestResponse(*serviceRequest),
		User: UserSummary{
			ID:        fmt.Sprintf("%d", user.ID),
			Name:      user.Name,
			Phone:     user.Phone,
			Role:      user.Role,
			Status:    user.Status,
			CreatedAt: user.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.UTC().Format(time.RFC3339),
		},
		Plan: PlanSummary{
			ID:          fmt.Sprintf("%d", plan.ID),
			ProviderID:  fmt.Sprintf("%d", plan.ProviderID),
			Name:        plan.Name,
			Description: plan.Description,
			Price:       plan.Price,
			Discount:    plan.Discount,
			IsActive:    plan.IsActive,
			CreatedAt:   plan.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:   plan.UpdatedAt.UTC().Format(time.RFC3339),
		},
		Provider: ProviderSummary{
			ID:        fmt.Sprintf("%d", provider.ID),
			Name:      provider.Name,
			ImageURL:  provider.ImageURL,
			CreatedAt: provider.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt: provider.UpdatedAt.UTC().Format(time.RFC3339),
		},
	}

	return response, nil
}

func (s *Service) UpdateStatus(ctx context.Context, id int64, targetStatus string) (*ServiceRequestResponse, error) {
	if id <= 0 {
		return nil, ErrServiceRequestInvalid
	}

	targetStatus = strings.TrimSpace(targetStatus)
	if !IsValidStatus(targetStatus) {
		return nil, ErrServiceRequestInvalid
	}

	var updated ServiceRequest
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var current ServiceRequest
		if err := tx.Where("id = ? AND deleted_at IS NULL", id).Take(&current).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrServiceRequestNotFound
			}
			return fmt.Errorf("get current service request: %w", err)
		}

		validTransition, transitionErr := CanTransition(current.Status, targetStatus)
		if transitionErr != nil {
			return ErrServiceRequestInvalid
		}
		if !validTransition {
			return ErrInvalidTransition
		}

		if current.Status != targetStatus {
			if err := tx.Model(&ServiceRequest{}).
				Where("id = ?", id).
				Updates(map[string]any{
					"status":     targetStatus,
					"updated_at": time.Now().UTC(),
				}).Error; err != nil {
				return fmt.Errorf("update service request status: %w", err)
			}

			if err := tx.Where("id = ?", id).Take(&updated).Error; err != nil {
				return fmt.Errorf("reload updated service request: %w", err)
			}
		} else {
			updated = current
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	response := toServiceRequestResponse(updated)
	return &response, nil
}

func (s *Service) repoGetByID(ctx context.Context, id int64) (*ServiceRequest, error) {
	var row ServiceRequest
	if err := s.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).Take(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrServiceRequestNotFound
		}
		return nil, fmt.Errorf("get service request: %w", err)
	}

	return &row, nil
}

func toServiceRequestResponse(row ServiceRequest) ServiceRequestResponse {
	return ServiceRequestResponse{
		ID:        fmt.Sprintf("%d", row.ID),
		UserID:    fmt.Sprintf("%d", row.UserID),
		PlanID:    fmt.Sprintf("%d", row.PlanID),
		Status:    row.Status,
		CreatedAt: row.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: row.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func normalizePaging(page int, pageSize int) (int, int, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 200 {
		return 0, 0, ErrServiceRequestInvalid
	}

	return page, pageSize, nil
}

func resolveSort(sort string) string {
	switch sort {
	case "created_at_asc":
		return "service_requests.created_at ASC"
	case "updated_at_desc":
		return "service_requests.updated_at DESC"
	case "updated_at_asc":
		return "service_requests.updated_at ASC"
	default:
		return "service_requests.created_at DESC"
	}
}
