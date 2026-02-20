package database

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/CallPilotReceptionist/internal/domain/entities"
	"github.com/CallPilotReceptionist/internal/domain/errors"
)

type InteractionRepositoryImpl struct {
	db *DB
}

func NewInteractionRepository(db *DB) InteractionRepository {
	return &InteractionRepositoryImpl{db: db}
}

func (r *InteractionRepositoryImpl) Create(ctx context.Context, interaction *entities.Interaction) error {
	interaction.ID = uuid.New().String()

	contentJSON, err := json.Marshal(interaction.Content)
	if err != nil {
		return errors.NewDatabaseError(err, "failed to marshal content")
	}

	query := `
		INSERT INTO interactions (id, call_id, type, content, timestamp, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = r.db.ExecContext(ctx, query,
		interaction.ID,
		interaction.CallID,
		interaction.Type,
		contentJSON,
		interaction.Timestamp,
		interaction.CreatedAt,
	)

	if err != nil {
		return errors.NewDatabaseError(err, "failed to create interaction")
	}

	return nil
}

func (r *InteractionRepositoryImpl) GetByID(ctx context.Context, id string) (*entities.Interaction, error) {
	query := `
		SELECT id, call_id, type, content, timestamp, created_at
		FROM interactions
		WHERE id = $1
	`

	interaction := &entities.Interaction{}
	var contentJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&interaction.ID,
		&interaction.CallID,
		&interaction.Type,
		&contentJSON,
		&interaction.Timestamp,
		&interaction.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.NewNotFoundError("interaction", id)
	}
	if err != nil {
		return nil, errors.NewDatabaseError(err, "failed to get interaction")
	}

	if err := json.Unmarshal(contentJSON, &interaction.Content); err != nil {
		return nil, errors.NewDatabaseError(err, "failed to unmarshal content")
	}

	return interaction, nil
}

func (r *InteractionRepositoryImpl) GetByCallID(ctx context.Context, callID string) ([]*entities.Interaction, error) {
	query := `
		SELECT id, call_id, type, content, timestamp, created_at
		FROM interactions
		WHERE call_id = $1
		ORDER BY timestamp ASC
	`

	rows, err := r.db.QueryContext(ctx, query, callID)
	if err != nil {
		return nil, errors.NewDatabaseError(err, "failed to get interactions by call")
	}
	defer rows.Close()

	return r.scanInteractions(rows)
}

func (r *InteractionRepositoryImpl) List(ctx context.Context, businessID string, limit, offset int) ([]*entities.Interaction, error) {
	query := `
		SELECT i.id, i.call_id, i.type, i.content, i.timestamp, i.created_at
		FROM interactions i
		JOIN calls c ON i.call_id = c.id
		WHERE c.business_id = $1
		ORDER BY i.timestamp DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, businessID, limit, offset)
	if err != nil {
		return nil, errors.NewDatabaseError(err, "failed to list interactions")
	}
	defer rows.Close()

	return r.scanInteractions(rows)
}

func (r *InteractionRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM interactions WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.NewDatabaseError(err, "failed to delete interaction")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewDatabaseError(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError("interaction", id)
	}

	return nil
}

func (r *InteractionRepositoryImpl) scanInteractions(rows *sql.Rows) ([]*entities.Interaction, error) {
	var interactions []*entities.Interaction

	for rows.Next() {
		interaction := &entities.Interaction{}
		var contentJSON []byte

		err := rows.Scan(
			&interaction.ID,
			&interaction.CallID,
			&interaction.Type,
			&contentJSON,
			&interaction.Timestamp,
			&interaction.CreatedAt,
		)
		if err != nil {
			return nil, errors.NewDatabaseError(err, "failed to scan interaction")
		}

		if err := json.Unmarshal(contentJSON, &interaction.Content); err != nil {
			return nil, errors.NewDatabaseError(err, "failed to unmarshal content")
		}

		interactions = append(interactions, interaction)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.NewDatabaseError(err, "failed to iterate interactions")
	}

	return interactions, nil
}
