package database

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/CallPilotReceptionist/internal/domain/entities"
	"github.com/CallPilotReceptionist/internal/domain/errors"
)

type BusinessRepositoryImpl struct {
	db *DB
}

func NewBusinessRepository(db *DB) BusinessRepository {
	return &BusinessRepositoryImpl{db: db}
}

func (r *BusinessRepositoryImpl) Create(ctx context.Context, business *entities.Business) error {
	business.ID = uuid.New().String()

	settingsJSON, err := json.Marshal(business.Settings)
	if err != nil {
		return errors.NewDatabaseError(err, "failed to marshal settings")
	}

	query := `
		INSERT INTO businesses (id, name, type, phone, settings, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = r.db.ExecContext(ctx, query,
		business.ID,
		business.Name,
		business.Type,
		business.Phone,
		settingsJSON,
		business.CreatedAt,
		business.UpdatedAt,
	)

	if err != nil {
		return errors.NewDatabaseError(err, "failed to create business")
	}

	return nil
}

func (r *BusinessRepositoryImpl) GetByID(ctx context.Context, id string) (*entities.Business, error) {
	query := `
		SELECT id, name, type, phone, settings, created_at, updated_at
		FROM businesses
		WHERE id = $1
	`

	business := &entities.Business{}
	var settingsJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&business.ID,
		&business.Name,
		&business.Type,
		&business.Phone,
		&settingsJSON,
		&business.CreatedAt,
		&business.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.NewNotFoundError("business", id)
	}
	if err != nil {
		return nil, errors.NewDatabaseError(err, "failed to get business")
	}

	if err := json.Unmarshal(settingsJSON, &business.Settings); err != nil {
		return nil, errors.NewDatabaseError(err, "failed to unmarshal settings")
	}

	return business, nil
}

func (r *BusinessRepositoryImpl) GetByPhone(ctx context.Context, phone string) (*entities.Business, error) {
	query := `
		SELECT id, name, type, phone, settings, created_at, updated_at
		FROM businesses
		WHERE phone = $1
	`

	business := &entities.Business{}
	var settingsJSON []byte

	err := r.db.QueryRowContext(ctx, query, phone).Scan(
		&business.ID,
		&business.Name,
		&business.Type,
		&business.Phone,
		&settingsJSON,
		&business.CreatedAt,
		&business.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.NewNotFoundError("business", phone)
	}
	if err != nil {
		return nil, errors.NewDatabaseError(err, "failed to get business")
	}

	if err := json.Unmarshal(settingsJSON, &business.Settings); err != nil {
		return nil, errors.NewDatabaseError(err, "failed to unmarshal settings")
	}

	return business, nil
}

func (r *BusinessRepositoryImpl) Update(ctx context.Context, business *entities.Business) error {
	settingsJSON, err := json.Marshal(business.Settings)
	if err != nil {
		return errors.NewDatabaseError(err, "failed to marshal settings")
	}

	query := `
		UPDATE businesses
		SET name = $2, type = $3, phone = $4, settings = $5, updated_at = $6
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		business.ID,
		business.Name,
		business.Type,
		business.Phone,
		settingsJSON,
		business.UpdatedAt,
	)

	if err != nil {
		return errors.NewDatabaseError(err, "failed to update business")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewDatabaseError(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError("business", business.ID)
	}

	return nil
}

func (r *BusinessRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM businesses WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.NewDatabaseError(err, "failed to delete business")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewDatabaseError(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError("business", id)
	}

	return nil
}

func (r *BusinessRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entities.Business, error) {
	query := `
		SELECT id, name, type, phone, settings, created_at, updated_at
		FROM businesses
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, errors.NewDatabaseError(err, "failed to list businesses")
	}
	defer rows.Close()

	var businesses []*entities.Business

	for rows.Next() {
		business := &entities.Business{}
		var settingsJSON []byte

		err := rows.Scan(
			&business.ID,
			&business.Name,
			&business.Type,
			&business.Phone,
			&settingsJSON,
			&business.CreatedAt,
			&business.UpdatedAt,
		)

		if err != nil {
			return nil, errors.NewDatabaseError(err, "failed to scan business")
		}

		if err := json.Unmarshal(settingsJSON, &business.Settings); err != nil {
			return nil, errors.NewDatabaseError(err, "failed to unmarshal settings")
		}

		businesses = append(businesses, business)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.NewDatabaseError(err, "failed to iterate businesses")
	}

	return businesses, nil
}
