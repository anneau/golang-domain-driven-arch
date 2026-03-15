package port

import (
	"context"

	"github.com/hkobori/golang-domain-driven-arch/internal/domain/entity"
	"github.com/hkobori/golang-domain-driven-arch/internal/domain/value"
)

type UserRepository interface {
	FindByID(ctx context.Context, id value.UserID) (*entity.User, error)
	FindByEmail(ctx context.Context, email value.Email) (*entity.User, error)
	Save(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id value.UserID) error
}
