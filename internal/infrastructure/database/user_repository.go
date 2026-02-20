package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/CallPilotReceptionist/internal/domain/entities"
	"github.com/CallPilotReceptionist/internal/domain/errors"
	"golang.org/x/crypto/bcrypt"
)

type UserRepositoryImpl struct {
	db *DB
}

func NewUserRepository(db *DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *entities.User) error {
	user.ID = uuid.New().String()

	query := `
		INSERT INTO users (id, business_id, email, password_hash, role, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.BusinessID,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.CreatedAt,
	)

	if err != nil {
		return errors.NewDatabaseError(err, "failed to create user")
	}

	return nil
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id string) (*entities.User, error) {
	query := `
		SELECT id, business_id, email, password_hash, role, created_at
		FROM users
		WHERE id = $1
	`

	user := &entities.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.BusinessID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.NewNotFoundError("user", id)
	}
	if err != nil {
		return nil, errors.NewDatabaseError(err, "failed to get user")
	}

	return user, nil
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := `
		SELECT id, business_id, email, password_hash, role, created_at
		FROM users
		WHERE email = $1
	`

	user := &entities.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.BusinessID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.NewNotFoundError("user", email)
	}
	if err != nil {
		return nil, errors.NewDatabaseError(err, "failed to get user by email")
	}

	return user, nil
}

func (r *UserRepositoryImpl) GetByBusinessID(ctx context.Context, businessID string) ([]*entities.User, error) {
	query := `
		SELECT id, business_id, email, password_hash, role, created_at
		FROM users
		WHERE business_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, businessID)
	if err != nil {
		return nil, errors.NewDatabaseError(err, "failed to get users by business")
	}
	defer rows.Close()

	var users []*entities.User

	for rows.Next() {
		user := &entities.User{}
		err := rows.Scan(
			&user.ID,
			&user.BusinessID,
			&user.Email,
			&user.PasswordHash,
			&user.Role,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, errors.NewDatabaseError(err, "failed to scan user")
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.NewDatabaseError(err, "failed to iterate users")
	}

	return users, nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, user *entities.User) error {
	query := `
		UPDATE users
		SET email = $2, password_hash = $3, role = $4
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Role,
	)

	if err != nil {
		return errors.NewDatabaseError(err, "failed to update user")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewDatabaseError(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError("user", user.ID)
	}

	return nil
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.NewDatabaseError(err, "failed to delete user")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewDatabaseError(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError("user", id)
	}

	return nil
}

// HashPassword hashes a plain text password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ComparePassword compares a hashed password with a plain text password
func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
