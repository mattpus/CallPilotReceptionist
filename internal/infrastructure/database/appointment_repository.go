package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/CallPilotReceptionist/internal/domain/entities"
	"github.com/CallPilotReceptionist/internal/domain/errors"
)

type AppointmentRepositoryImpl struct {
	db *DB
}

func NewAppointmentRepository(db *DB) AppointmentRepository {
	return &AppointmentRepositoryImpl{db: db}
}

func (r *AppointmentRepositoryImpl) Create(ctx context.Context, appointment *entities.AppointmentRequest) error {
	appointment.ID = uuid.New().String()

	query := `
		INSERT INTO appointments (id, call_id, business_id, customer_name, customer_phone, 
			requested_date, requested_time, service_type, notes, status, extracted_at, confirmed_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := r.db.ExecContext(ctx, query,
		appointment.ID,
		appointment.CallID,
		appointment.BusinessID,
		appointment.CustomerName,
		appointment.CustomerPhone,
		appointment.RequestedDate,
		appointment.RequestedTime,
		appointment.ServiceType,
		appointment.Notes,
		appointment.Status,
		appointment.ExtractedAt,
		appointment.ConfirmedAt,
		appointment.CreatedAt,
	)

	if err != nil {
		return errors.NewDatabaseError(err, "failed to create appointment")
	}

	return nil
}

func (r *AppointmentRepositoryImpl) GetByID(ctx context.Context, id string) (*entities.AppointmentRequest, error) {
	query := `
		SELECT id, call_id, business_id, customer_name, customer_phone, 
			requested_date, requested_time, service_type, notes, status, extracted_at, confirmed_at, created_at
		FROM appointments
		WHERE id = $1
	`

	appointment := &entities.AppointmentRequest{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&appointment.ID,
		&appointment.CallID,
		&appointment.BusinessID,
		&appointment.CustomerName,
		&appointment.CustomerPhone,
		&appointment.RequestedDate,
		&appointment.RequestedTime,
		&appointment.ServiceType,
		&appointment.Notes,
		&appointment.Status,
		&appointment.ExtractedAt,
		&appointment.ConfirmedAt,
		&appointment.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.NewNotFoundError("appointment", id)
	}
	if err != nil {
		return nil, errors.NewDatabaseError(err, "failed to get appointment")
	}

	return appointment, nil
}

func (r *AppointmentRepositoryImpl) GetByCallID(ctx context.Context, callID string) (*entities.AppointmentRequest, error) {
	query := `
		SELECT id, call_id, business_id, customer_name, customer_phone, 
			requested_date, requested_time, service_type, notes, status, extracted_at, confirmed_at, created_at
		FROM appointments
		WHERE call_id = $1
	`

	appointment := &entities.AppointmentRequest{}
	err := r.db.QueryRowContext(ctx, query, callID).Scan(
		&appointment.ID,
		&appointment.CallID,
		&appointment.BusinessID,
		&appointment.CustomerName,
		&appointment.CustomerPhone,
		&appointment.RequestedDate,
		&appointment.RequestedTime,
		&appointment.ServiceType,
		&appointment.Notes,
		&appointment.Status,
		&appointment.ExtractedAt,
		&appointment.ConfirmedAt,
		&appointment.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.NewNotFoundError("appointment", callID)
	}
	if err != nil {
		return nil, errors.NewDatabaseError(err, "failed to get appointment by call")
	}

	return appointment, nil
}

func (r *AppointmentRepositoryImpl) GetByBusinessID(ctx context.Context, businessID string, limit, offset int) ([]*entities.AppointmentRequest, error) {
	query := `
		SELECT id, call_id, business_id, customer_name, customer_phone, 
			requested_date, requested_time, service_type, notes, status, extracted_at, confirmed_at, created_at
		FROM appointments
		WHERE business_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, businessID, limit, offset)
	if err != nil {
		return nil, errors.NewDatabaseError(err, "failed to get appointments by business")
	}
	defer rows.Close()

	return r.scanAppointments(rows)
}

func (r *AppointmentRepositoryImpl) GetPendingAppointments(ctx context.Context, businessID string) ([]*entities.AppointmentRequest, error) {
	query := `
		SELECT id, call_id, business_id, customer_name, customer_phone, 
			requested_date, requested_time, service_type, notes, status, extracted_at, confirmed_at, created_at
		FROM appointments
		WHERE business_id = $1 AND status = 'pending'
		ORDER BY extracted_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, businessID)
	if err != nil {
		return nil, errors.NewDatabaseError(err, "failed to get pending appointments")
	}
	defer rows.Close()

	return r.scanAppointments(rows)
}

func (r *AppointmentRepositoryImpl) Update(ctx context.Context, appointment *entities.AppointmentRequest) error {
	query := `
		UPDATE appointments
		SET customer_name = $2, customer_phone = $3, requested_date = $4, requested_time = $5, 
			service_type = $6, notes = $7, status = $8, confirmed_at = $9
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		appointment.ID,
		appointment.CustomerName,
		appointment.CustomerPhone,
		appointment.RequestedDate,
		appointment.RequestedTime,
		appointment.ServiceType,
		appointment.Notes,
		appointment.Status,
		appointment.ConfirmedAt,
	)

	if err != nil {
		return errors.NewDatabaseError(err, "failed to update appointment")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewDatabaseError(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError("appointment", appointment.ID)
	}

	return nil
}

func (r *AppointmentRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM appointments WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.NewDatabaseError(err, "failed to delete appointment")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewDatabaseError(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError("appointment", id)
	}

	return nil
}

func (r *AppointmentRepositoryImpl) scanAppointments(rows *sql.Rows) ([]*entities.AppointmentRequest, error) {
	var appointments []*entities.AppointmentRequest

	for rows.Next() {
		appointment := &entities.AppointmentRequest{}
		err := rows.Scan(
			&appointment.ID,
			&appointment.CallID,
			&appointment.BusinessID,
			&appointment.CustomerName,
			&appointment.CustomerPhone,
			&appointment.RequestedDate,
			&appointment.RequestedTime,
			&appointment.ServiceType,
			&appointment.Notes,
			&appointment.Status,
			&appointment.ExtractedAt,
			&appointment.ConfirmedAt,
			&appointment.CreatedAt,
		)
		if err != nil {
			return nil, errors.NewDatabaseError(err, "failed to scan appointment")
		}
		appointments = append(appointments, appointment)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.NewDatabaseError(err, "failed to iterate appointments")
	}

	return appointments, nil
}
