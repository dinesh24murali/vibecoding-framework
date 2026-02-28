package plans

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrPlanNotFound      = errors.New("plan not found")
	ErrPlanDuplicate     = errors.New("plan duplicate")
	ErrPlanValidation    = errors.New("plan validation error")
	ErrProviderNotFound  = errors.New("provider not found")
	ErrPlanLimitExceeded = errors.New("plan limit exceeded")
)

type Service struct {
	repo *Repository
}

type CreateInput struct {
	ProviderID  int64
	Name        string
	Description string
	Price       float64
	Discount    float64
	IsActive    bool
}

type UpdateInput struct {
	ProviderID  *int64
	Name        *string
	Description *string
	Price       *float64
	Discount    *float64
	IsActive    *bool
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListAdmin(ctx context.Context, page int, pageSize int, providerID *int64, isActive *bool, sort string) ([]Plan, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	return s.repo.ListAdmin(ctx, page, pageSize, providerID, isActive, sort)
}

func (s *Service) ListCustomerByProvider(ctx context.Context, providerID int64) ([]Plan, error) {
	exists, err := s.repo.ProviderExists(ctx, providerID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrProviderNotFound
	}

	return s.repo.ListCustomerByProvider(ctx, providerID)
}

func (s *Service) GetByID(ctx context.Context, id int64) (*Plan, error) {
	plan, err := s.repo.FindActiveByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, ErrPlanNotFound
	}

	return plan, nil
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*Plan, error) {
	if err := validateCreateInput(input); err != nil {
		return nil, err
	}

	exists, err := s.repo.ProviderExists(ctx, input.ProviderID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrProviderNotFound
	}

	count, err := s.repo.ActivePlanCountByProvider(ctx, input.ProviderID, nil)
	if err != nil {
		return nil, err
	}
	if count >= 100 {
		return nil, ErrPlanLimitExceeded
	}

	duplicate, err := s.repo.ExistsActiveByProviderAndName(ctx, input.ProviderID, input.Name, nil)
	if err != nil {
		return nil, err
	}
	if duplicate {
		return nil, ErrPlanDuplicate
	}

	plan := Plan{
		ProviderID:  input.ProviderID,
		Name:        strings.TrimSpace(input.Name),
		Description: strings.TrimSpace(input.Description),
		Price:       input.Price,
		Discount:    input.Discount,
		IsActive:    input.IsActive,
	}

	if err := gorm.G[Plan](s.repo.Db.WithContext(ctx)).Create(ctx, &plan); err != nil {
		return nil, fmt.Errorf("create plan: %w", err)
	}

	return &plan, nil
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (*Plan, error) {
	plan, err := s.repo.FindActiveByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, ErrPlanNotFound
	}

	if input.ProviderID != nil {
		exists, existsErr := s.repo.ProviderExists(ctx, *input.ProviderID)
		if existsErr != nil {
			return nil, existsErr
		}
		if !exists {
			return nil, ErrProviderNotFound
		}
		plan.ProviderID = *input.ProviderID
	}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" || len(name) > 120 {
			return nil, ErrPlanValidation
		}
		plan.Name = name
	}

	if input.Description != nil {
		description := strings.TrimSpace(*input.Description)
		if description == "" || len(description) > 1000 {
			return nil, ErrPlanValidation
		}
		plan.Description = description
	}

	if input.Price != nil {
		if *input.Price < 0.01 {
			return nil, ErrPlanValidation
		}
		plan.Price = *input.Price
	}

	if input.Discount != nil {
		if *input.Discount < 0 {
			return nil, ErrPlanValidation
		}
		plan.Discount = *input.Discount
	}

	if plan.Discount > plan.Price {
		return nil, ErrPlanValidation
	}

	if input.IsActive != nil {
		plan.IsActive = *input.IsActive
	}

	excludeID := plan.ID
	duplicate, err := s.repo.ExistsActiveByProviderAndName(ctx, plan.ProviderID, plan.Name, &excludeID)
	if err != nil {
		return nil, err
	}
	if duplicate {
		return nil, ErrPlanDuplicate
	}

	count, err := s.repo.ActivePlanCountByProvider(ctx, plan.ProviderID, &excludeID)
	if err != nil {
		return nil, err
	}
	if count >= 100 {
		return nil, ErrPlanLimitExceeded
	}

	if err := s.repo.Db.WithContext(ctx).Save(plan).Error; err != nil {
		return nil, fmt.Errorf("update plan: %w", err)
	}

	return plan, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	deleted, err := s.repo.SoftDelete(ctx, id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrPlanNotFound
	}

	return nil
}

func validateCreateInput(input CreateInput) error {
	if input.ProviderID <= 0 {
		return ErrPlanValidation
	}

	if strings.TrimSpace(input.Name) == "" || len(strings.TrimSpace(input.Name)) > 120 {
		return ErrPlanValidation
	}

	if strings.TrimSpace(input.Description) == "" || len(strings.TrimSpace(input.Description)) > 1000 {
		return ErrPlanValidation
	}

	if input.Price < 0.01 {
		return ErrPlanValidation
	}

	if input.Discount < 0 || input.Discount > input.Price {
		return ErrPlanValidation
	}

	return nil
}
