package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/CallPilotReceptionist/internal/domain/entities"
	"github.com/CallPilotReceptionist/internal/domain/errors"
)

type TranscriptRepositoryImpl struct {
	db *DB
}

func NewTranscriptRepository(db *DB) TranscriptRepository {
	return &TranscriptRepositoryImpl{db: db}
}

func (r *TranscriptRepositoryImpl) Create(ctx context.Context, transcript *entities.Transcript) error {
	transcript.ID = uuid.New().String()

	query := `
		INSERT INTO transcripts (id, call_id, role, message, timestamp, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		transcript.ID,
		transcript.CallID,
		transcript.Role,
		transcript.Message,
		transcript.Timestamp,
		transcript.CreatedAt,
	)

	if err != nil {
		return errors.NewDatabaseError(err, "failed to create transcript")
	}

	return nil
}

func (r *TranscriptRepositoryImpl) CreateBatch(ctx context.Context, transcripts []*entities.Transcript) error {
	if len(transcripts) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.NewDatabaseError(err, "failed to begin transaction")
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO transcripts (id, call_id, role, message, timestamp, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`)
	if err != nil {
		return errors.NewDatabaseError(err, "failed to prepare statement")
	}
	defer stmt.Close()

	for _, transcript := range transcripts {
		transcript.ID = uuid.New().String()
		_, err := stmt.ExecContext(ctx,
			transcript.ID,
			transcript.CallID,
			transcript.Role,
			transcript.Message,
			transcript.Timestamp,
			transcript.CreatedAt,
		)
		if err != nil {
			return errors.NewDatabaseError(err, "failed to insert transcript in batch")
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.NewDatabaseError(err, "failed to commit transaction")
	}

	return nil
}

func (r *TranscriptRepositoryImpl) GetByCallID(ctx context.Context, callID string) ([]*entities.Transcript, error) {
	query := `
		SELECT id, call_id, role, message, timestamp, created_at
		FROM transcripts
		WHERE call_id = $1
		ORDER BY timestamp ASC
	`

	rows, err := r.db.QueryContext(ctx, query, callID)
	if err != nil {
		return nil, errors.NewDatabaseError(err, "failed to get transcripts by call")
	}
	defer rows.Close()

	var transcripts []*entities.Transcript

	for rows.Next() {
		transcript := &entities.Transcript{}
		err := rows.Scan(
			&transcript.ID,
			&transcript.CallID,
			&transcript.Role,
			&transcript.Message,
			&transcript.Timestamp,
			&transcript.CreatedAt,
		)
		if err != nil {
			return nil, errors.NewDatabaseError(err, "failed to scan transcript")
		}
		transcripts = append(transcripts, transcript)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.NewDatabaseError(err, "failed to iterate transcripts")
	}

	return transcripts, nil
}

func (r *TranscriptRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM transcripts WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.NewDatabaseError(err, "failed to delete transcript")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewDatabaseError(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError("transcript", id)
	}

	return nil
}

func (r *TranscriptRepositoryImpl) DeleteByCallID(ctx context.Context, callID string) error {
	query := `DELETE FROM transcripts WHERE call_id = $1`

	_, err := r.db.ExecContext(ctx, query, callID)
	if err != nil {
		return errors.NewDatabaseError(err, "failed to delete transcripts by call")
	}

	return nil
}
