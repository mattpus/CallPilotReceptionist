package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/CallPilotReceptionist/internal/domain/entities"
	"github.com/CallPilotReceptionist/internal/domain/errors"
)

type CallRepositoryImpl struct {
	db *DB
}

func NewCallRepository(db *DB) CallRepository {
	return &CallRepositoryImpl{db: db}
}

func (r *CallRepositoryImpl) Create(ctx context.Context, call *entities.Call) error {
	call.ID = uuid.New().String()

	query := `
		INSERT INTO calls (id, business_id, provider_call_id, caller_phone, duration, status, cost, started_at, ended_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.ExecContext(ctx, query,
		call.ID,
		call.BusinessID,
		call.ProviderCallID,
		call.CallerPhone,
		call.Duration,
		call.Status,
		call.Cost,
		call.StartedAt,
		call.EndedAt,
		call.CreatedAt,
	)

	if err != nil {
		return errors.NewDatabaseError(err, "failed to create call")
	}

	return nil
}

func (r *CallRepositoryImpl) GetByID(ctx context.Context, id string) (*entities.Call, error) {
	query := `
		SELECT id, business_id, provider_call_id, caller_phone, duration, status, cost, started_at, ended_at, created_at
		FROM calls
		WHERE id = $1
	`

	call := &entities.Call{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&call.ID,
		&call.BusinessID,
		&call.ProviderCallID,
		&call.CallerPhone,
		&call.Duration,
		&call.Status,
		&call.Cost,
		&call.StartedAt,
		&call.EndedAt,
		&call.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.NewNotFoundError("call", id)
	}
	if err != nil {
		return nil, errors.NewDatabaseError(err, "failed to get call")
	}

	return call, nil
}

func (r *CallRepositoryImpl) GetByProviderCallID(ctx context.Context, providerCallID string) (*entities.Call, error) {
	query := `
		SELECT id, business_id, provider_call_id, caller_phone, duration, status, cost, started_at, ended_at, created_at
		FROM calls
		WHERE provider_call_id = $1
	`

	call := &entities.Call{}
	err := r.db.QueryRowContext(ctx, query, providerCallID).Scan(
		&call.ID,
		&call.BusinessID,
		&call.ProviderCallID,
		&call.CallerPhone,
		&call.Duration,
		&call.Status,
		&call.Cost,
		&call.StartedAt,
		&call.EndedAt,
		&call.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.NewNotFoundError("call", providerCallID)
	}
	if err != nil {
		return nil, errors.NewDatabaseError(err, "failed to get call by provider ID")
	}

	return call, nil
}

func (r *CallRepositoryImpl) GetByBusinessID(ctx context.Context, businessID string, limit, offset int) ([]*entities.Call, error) {
	query := `
		SELECT id, business_id, provider_call_id, caller_phone, duration, status, cost, started_at, ended_at, created_at
		FROM calls
		WHERE business_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, businessID, limit, offset)
	if err != nil {
		return nil, errors.NewDatabaseError(err, "failed to get calls by business")
	}
	defer rows.Close()

	return r.scanCalls(rows)
}

func (r *CallRepositoryImpl) Update(ctx context.Context, call *entities.Call) error {
	query := `
		UPDATE calls
		SET provider_call_id = $2, caller_phone = $3, duration = $4, status = $5, cost = $6, started_at = $7, ended_at = $8
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		call.ID,
		call.ProviderCallID,
		call.CallerPhone,
		call.Duration,
		call.Status,
		call.Cost,
		call.StartedAt,
		call.EndedAt,
	)

	if err != nil {
		return errors.NewDatabaseError(err, "failed to update call")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewDatabaseError(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError("call", call.ID)
	}

	return nil
}

func (r *CallRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM calls WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.NewDatabaseError(err, "failed to delete call")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewDatabaseError(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError("call", id)
	}

	return nil
}

func (r *CallRepositoryImpl) GetByDateRange(ctx context.Context, businessID string, startDate, endDate time.Time) ([]*entities.Call, error) {
	query := `
		SELECT id, business_id, provider_call_id, caller_phone, duration, status, cost, started_at, ended_at, created_at
		FROM calls
		WHERE business_id = $1 AND created_at BETWEEN $2 AND $3
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, businessID, startDate, endDate)
	if err != nil {
		return nil, errors.NewDatabaseError(err, "failed to get calls by date range")
	}
	defer rows.Close()

	return r.scanCalls(rows)
}

func (r *CallRepositoryImpl) GetStats(ctx context.Context, businessID string, startDate, endDate time.Time) (*CallStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_calls,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_calls,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_calls,
			COALESCE(SUM(duration), 0) as total_duration,
			COALESCE(AVG(duration), 0) as average_duration,
			COALESCE(SUM(cost), 0) as total_cost
		FROM calls
		WHERE business_id = $1 AND created_at BETWEEN $2 AND $3
	`

	stats := &CallStats{}
	err := r.db.QueryRowContext(ctx, query, businessID, startDate, endDate).Scan(
		&stats.TotalCalls,
		&stats.CompletedCalls,
		&stats.FailedCalls,
		&stats.TotalDuration,
		&stats.AverageDuration,
		&stats.TotalCost,
	)

	if err != nil {
		return nil, errors.NewDatabaseError(err, "failed to get call stats")
	}

	return stats, nil
}

func (r *CallRepositoryImpl) scanCalls(rows *sql.Rows) ([]*entities.Call, error) {
	var calls []*entities.Call

	for rows.Next() {
		call := &entities.Call{}
		err := rows.Scan(
			&call.ID,
			&call.BusinessID,
			&call.ProviderCallID,
			&call.CallerPhone,
			&call.Duration,
			&call.Status,
			&call.Cost,
			&call.StartedAt,
			&call.EndedAt,
			&call.CreatedAt,
		)
		if err != nil {
			return nil, errors.NewDatabaseError(err, "failed to scan call")
		}
		calls = append(calls, call)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.NewDatabaseError(err, "failed to iterate calls")
	}

	return calls, nil
}
