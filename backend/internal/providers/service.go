package providers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrProviderNotFound   = errors.New("provider not found")
	ErrProviderDuplicate  = errors.New("provider duplicate")
	ErrProviderValidation = errors.New("provider validation error")
)

type Service struct {
	repo *Repository
}

type CreateInput struct {
	Name     string
	ImageURL *string
}

type UpdateInput struct {
	Name     *string
	ImageURL *string
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListAdmin(ctx context.Context, page int, pageSize int, search string, sort string) ([]Provider, int64, error) {
	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	return s.repo.ListAdmin(ctx, page, pageSize, search, sort)
}

func (s *Service) ListCustomer(ctx context.Context) ([]Provider, error) {
	return s.repo.ListCustomerActive(ctx)
}

func (s *Service) GetByID(ctx context.Context, id int64) (*Provider, error) {
	provider, err := s.repo.FindActiveByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if provider == nil {
		return nil, ErrProviderNotFound
	}

	return provider, nil
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*Provider, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" || len(name) > 120 {
		return nil, ErrProviderValidation
	}

	exists, err := s.repo.ExistsActiveByName(ctx, name, nil)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, ErrProviderDuplicate
	}

	provider := Provider{Name: name, ImageURL: input.ImageURL}
	if err := gorm.G[Provider](s.repo.Db.WithContext(ctx)).Create(ctx, &provider); err != nil {
		return nil, fmt.Errorf("create provider: %w", err)
	}

	return &provider, nil
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (*Provider, error) {
	provider, err := s.repo.FindActiveByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if provider == nil {
		return nil, ErrProviderNotFound
	}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" || len(name) > 120 {
			return nil, ErrProviderValidation
		}

		excludeID := provider.ID
		exists, existsErr := s.repo.ExistsActiveByName(ctx, name, &excludeID)
		if existsErr != nil {
			return nil, existsErr
		}

		if exists {
			return nil, ErrProviderDuplicate
		}

		provider.Name = name
	}

	if input.ImageURL != nil {
		provider.ImageURL = input.ImageURL
	}

	if err := s.repo.Db.WithContext(ctx).Save(provider).Error; err != nil {
		return nil, fmt.Errorf("update provider: %w", err)
	}

	return provider, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	deleted, err := s.repo.SoftDelete(ctx, id)
	if err != nil {
		return err
	}

	if !deleted {
		return ErrProviderNotFound
	}

	return nil
}
