package repository

import (
	"context"
	"errors"

	"github.com/alexnakagama/server-monitor/internal/model"
	"github.com/alexnakagama/server-monitor/internal/server/errors_custom"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, user model.User) error {
	query := `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
	)

	return err
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (model.User, error) {
	var user model.User

	query := `
		SELECT id, username, email, password_hash, created_at
		FROM users
		WHERE username = $1
	`

	err := r.db.QueryRow(
		ctx,
		query,
		username,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return model.User{}, errors_custom.ErrUserNotFound
	}

	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User

	query := `
		SELECT id, username, email, password_hash, created_at
		FROM users
		WHERE email = $1
	`

	err := r.db.QueryRow(
		ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return model.User{}, errors_custom.ErrUserNotFound
	}

	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors_custom.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) Update(ctx context.Context, user model.User) error {
	query := `
		UPDATE users
		SET username = $1,
		email = $2
		WHERE id = $3
	`

	result, err := r.db.Exec(
		ctx,
		query,
		user.Username,
		user.Email,
		user.ID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors_custom.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (model.User, error) {
	var user model.User

	query := `
		SELECT id, username, email, password_hash, created_at
		FROM users
		WHERE id = $1
	`

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return model.User{}, ErrUserNotFound
	}

	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
