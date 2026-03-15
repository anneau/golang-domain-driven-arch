package database

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type rawPublisher interface {
	PublishRaw(ctx context.Context, payload []byte) error
}

type OutboxRelay struct {
	db        *sql.DB
	publisher rawPublisher
	interval  time.Duration
}

func NewOutboxRelay(db *sql.DB, publisher rawPublisher, interval time.Duration) *OutboxRelay {
	return &OutboxRelay{db: db, publisher: publisher, interval: interval}
}

func (r *OutboxRelay) Run(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := r.relay(ctx); err != nil {
				log.Printf("outbox relay error: %v", err)
			}
		}
	}
}

func (r *OutboxRelay) relay(ctx context.Context) error {
	records, err := findUnpublishedOutbox(ctx, r.db, 100)
	if err != nil {
		return err
	}

	for _, rec := range records {
		if err := r.publisher.PublishRaw(ctx, rec.payload); err != nil {
			return err
		}
		if err := markOutboxPublished(ctx, r.db, rec.id); err != nil {
			return err
		}
	}
	return nil
}
