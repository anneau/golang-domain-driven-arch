package database

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	messagingevent "github.com/hkobori/golang-domain-driven-arch/internal/adapter/messaging/event"
	"github.com/hkobori/golang-domain-driven-arch/internal/domain/entity"
	"github.com/hkobori/golang-domain-driven-arch/internal/domain/event"
	"github.com/hkobori/golang-domain-driven-arch/internal/domain/value"
)

type userRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepositoryImpl {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) FindByID(ctx context.Context, id value.UserID) (*entity.User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, name, email FROM users WHERE id = $1`,
		id.String(),
	)
	return r.scan(row)
}

func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email value.Email) (*entity.User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, name, email FROM users WHERE email = $1`,
		email.String(),
	)
	return r.scan(row)
}

func (r *userRepositoryImpl) Save(ctx context.Context, user *entity.User) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`INSERT INTO users (id, name, email)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, email = EXCLUDED.email`,
		user.ID().String(),
		user.Name(),
		user.Email().String(),
	)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	for _, evt := range user.UncommittedEvents() {
		if err := insertOutbox(ctx, tx, evt); err != nil {
			return fmt.Errorf("failed to save outbox event: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (r *userRepositoryImpl) Delete(ctx context.Context, id value.UserID) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM users WHERE id = $1`,
		id.String(),
	)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (r *userRepositoryImpl) scan(row *sql.Row) (*entity.User, error) {
	var id, name, email string
	if err := row.Scan(&id, &name, &email); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan user row: %w", err)
	}

	userID, err := value.UserIDFrom(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user id in db: %w", err)
	}

	return entity.ReconstructUser(userID, name, value.EmailFrom(email)), nil
}

func insertOutbox(ctx context.Context, tx *sql.Tx, evt event.Event) error {
	payload, err := serializeEvent(evt)
	if err != nil {
		return err
	}

	id, err := newOutboxID()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO outbox (id, event_type, aggregate_id, payload, occurred_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		id,
		evt.EventType(),
		evt.AggregateID(),
		payload,
		evt.OccurredAt(),
	)
	return err
}

func serializeEvent(evt event.Event) ([]byte, error) {
	switch e := evt.(type) {
	case *event.UserCreatedEvent:
		return json.Marshal(messagingevent.UserCreatedEventDTO{
			EventType:  e.EventType(),
			UserID:     e.UserID(),
			Name:       e.Name(),
			Email:      e.Email(),
			OccurredAt: e.OccurredAt(),
		})
	default:
		return nil, fmt.Errorf("unknown event type: %T", evt)
	}
}

func newOutboxID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate outbox id: %w", err)
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

// outboxRecord is used internally by the relay.
type outboxRecord struct {
	id      string
	payload []byte
}

func findUnpublishedOutbox(ctx context.Context, db *sql.DB, limit int) ([]outboxRecord, error) {
	rows, err := db.QueryContext(ctx,
		`SELECT id, payload FROM outbox WHERE published_at IS NULL ORDER BY occurred_at LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []outboxRecord
	for rows.Next() {
		var r outboxRecord
		if err := rows.Scan(&r.id, &r.payload); err != nil {
			return nil, err
		}
		records = append(records, r)
	}
	return records, rows.Err()
}

func markOutboxPublished(ctx context.Context, db *sql.DB, id string) error {
	_, err := db.ExecContext(ctx,
		`UPDATE outbox SET published_at = $1 WHERE id = $2`,
		time.Now(),
		id,
	)
	return err
}
