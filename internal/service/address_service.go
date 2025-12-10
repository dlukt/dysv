package service

import (
	"context"
	"time"

	"github.com/deicod/dysv/internal/model"
	"github.com/deicod/dysv/internal/repo"
)

type AddressService struct {
	repo repo.AddressRepository
}

func NewAddressService(repo repo.AddressRepository) *AddressService {
	return &AddressService{repo: repo}
}

func (s *AddressService) CreateAddress(ctx context.Context, addr *model.Address) error {
	addr.CreatedAt = time.Now()
	addr.UpdatedAt = time.Now()

	// Business Logic: If default, unset others first
	if addr.IsDefault {
		if err := s.repo.UnsetDefaults(ctx, addr.UserID); err != nil {
			return err
		}
	}

	// Validation logic could go here (e.g. check country code)

	return s.repo.Create(ctx, addr)
}

func (s *AddressService) UpdateAddress(ctx context.Context, addr *model.Address) error {
	addr.UpdatedAt = time.Now()

	if addr.IsDefault {
		if err := s.repo.UnsetDefaults(ctx, addr.UserID); err != nil {
			return err
		}
	}

	return s.repo.Update(ctx, addr)
}

func (s *AddressService) ListAddresses(ctx context.Context, userID string) ([]model.Address, error) {
	return s.repo.ListByUserID(ctx, userID)
}

func (s *AddressService) DeleteAddress(ctx context.Context, id, userID string) error {
	return s.repo.Delete(ctx, id, userID)
}
