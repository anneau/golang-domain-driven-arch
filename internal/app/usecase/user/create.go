package user_usecase

import (
	"context"
	"fmt"

	"github.com/hkobori/golang-domain-driven-arch/internal/domain/entity"
	"github.com/hkobori/golang-domain-driven-arch/internal/domain/port"
	"github.com/hkobori/golang-domain-driven-arch/internal/domain/service"
	"github.com/hkobori/golang-domain-driven-arch/internal/domain/value"
)

type CreateUserUseCase interface {
	Execute(ctx context.Context, input CreateUserInput) (*UserOutput, error)
}

type createUserUseCaseImpl struct {
	userRepo       port.UserRepository
	userService    service.UserDomainServiceInterface
	eventPublisher port.EventPublisher
}

func NewCreateUserUseCase(
	userRepo port.UserRepository,
	userService service.UserDomainServiceInterface,
	eventPublisher port.EventPublisher,
) CreateUserUseCase {
	return &createUserUseCaseImpl{
		userRepo:       userRepo,
		userService:    userService,
		eventPublisher: eventPublisher,
	}
}

func (u *createUserUseCaseImpl) Execute(ctx context.Context, input CreateUserInput) (*UserOutput, error) {
	email, err := value.NewEmail(input.Email)
	if err != nil {
		return nil, &CreateUserError{Kind: ErrValidation, Message: err.Error()}
	}

	duplicated, err := u.userService.IsEmailDuplicated(ctx, email)
	if err != nil {
		return nil, err
	}
	if duplicated {
		return nil, &CreateUserError{Kind: ErrEmailDuplicated, Message: fmt.Sprintf("email %s is already registered", input.Email)}
	}

	id, err := value.NewUserID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate user id: %w", err)
	}

	userEntity, err := entity.NewUser(id, input.Name, email)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if err := u.userRepo.Save(ctx, userEntity); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	for _, evt := range userEntity.UncommittedEvents() {
		if err := u.eventPublisher.Publish(ctx, evt); err != nil {
			return nil, fmt.Errorf("failed to publish user created event: %w", err)
		}
	}
	userEntity.ClearEvents()

	return toOutput(userEntity), nil
}
