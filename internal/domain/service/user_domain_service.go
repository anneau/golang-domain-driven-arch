package service

import (
	"context"
	"fmt"

	"github.com/hkobori/golang-domain-driven-arch/internal/domain/port"
	"github.com/hkobori/golang-domain-driven-arch/internal/domain/value"
)

type UserDomainServiceInterface interface {
	IsEmailDuplicated(ctx context.Context, email value.Email) (bool, error)
	IsEmailDuplicatedExcluding(ctx context.Context, email value.Email, excludeID value.UserID) (bool, error)
}

type UserDomainService struct {
	userRepo port.UserRepository
}

func NewUserDomainService(userRepo port.UserRepository) *UserDomainService {
	return &UserDomainService{userRepo: userRepo}
}

func (s *UserDomainService) IsEmailDuplicated(ctx context.Context, email value.Email) (bool, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return false, fmt.Errorf("failed to check duplicate email: %w", err)
	}
	return user != nil, nil
}

func (s *UserDomainService) IsEmailDuplicatedExcluding(ctx context.Context, email value.Email, excludeID value.UserID) (bool, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return false, fmt.Errorf("failed to check duplicate email: %w", err)
	}
	if user == nil {
		return false, nil
	}
	return !user.ID().Equals(excludeID), nil
}

