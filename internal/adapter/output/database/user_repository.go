package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hkobori/golang-domain-driven-arch/internal/domain/entity"
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
	_, err := r.db.ExecContext(ctx,
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
