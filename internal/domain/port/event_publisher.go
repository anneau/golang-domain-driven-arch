package port

import (
	"context"

	"github.com/hkobori/golang-domain-driven-arch/internal/domain/event"
)

type EventPublisher interface {
	Publish(ctx context.Context, evt event.Event) error
}
