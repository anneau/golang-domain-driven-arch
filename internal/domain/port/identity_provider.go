package port

import (
	"context"

	"github.com/hkobori/golang-domain-driven-arch/internal/domain/entity"
)

type IdentityProvider interface {
	Save(ctx context.Context, identity *entity.Identity) error
}
