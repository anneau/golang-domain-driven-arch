package user_usecase

import (
	"github.com/hkobori/golang-domain-driven-arch/internal/domain/entity"
	"github.com/hkobori/golang-domain-driven-arch/internal/domain/port"
	"github.com/hkobori/golang-domain-driven-arch/internal/domain/service"
)

type CreateUserInput struct {
	Name  string
	Email string
}

type UpdateUserInput struct {
	ID    string
	Name  string
	Email string
}

type UserOutput struct {
	ID    string
	Name  string
	Email string
}

func toOutput(u *entity.User) *UserOutput {
	return &UserOutput{
		ID:    u.ID().String(),
		Name:  u.Name(),
		Email: u.Email().String(),
	}
}

type UserUseCases struct {
	Create CreateUserUseCase
}

func NewUserUseCases(
	userRepo port.UserRepository,
	userService service.UserDomainServiceInterface,
	eventPublisher port.EventPublisher,
) *UserUseCases {
	return &UserUseCases{
		Create: NewCreateUserUseCase(userRepo, userService, eventPublisher),
	}
}
