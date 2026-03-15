package identity_usecase

import (
	"context"
	"fmt"

	"github.com/hkobori/golang-domain-driven-arch/internal/domain/entity"
	"github.com/hkobori/golang-domain-driven-arch/internal/domain/port"
	"github.com/hkobori/golang-domain-driven-arch/internal/domain/value"
)

type RegisterIdentityInput struct {
	UserID string
	Email  string
}

type RegisterIdentityUseCase interface {
	Execute(ctx context.Context, input RegisterIdentityInput) error
}

type registerIdentityUseCaseImpl struct {
	identityProvider port.IdentityProvider
}

func NewRegisterIdentityUseCase(identityProvider port.IdentityProvider) RegisterIdentityUseCase {
	return &registerIdentityUseCaseImpl{
		identityProvider: identityProvider,
	}
}

func (u *registerIdentityUseCaseImpl) Execute(ctx context.Context, input RegisterIdentityInput) error {
	userID, err := value.UserIDFrom(input.UserID)
	if err != nil {
		return fmt.Errorf("failed to parse user id: %w", err)
	}

	email, err := value.NewEmail(input.Email)
	if err != nil {
		return fmt.Errorf("failed to parse email: %w", err)
	}

	identity := entity.NewIdentity(userID, email)

	if err := u.identityProvider.Save(ctx, identity); err != nil {
		return fmt.Errorf("failed to save identity: %w", err)
	}

	return nil
}
