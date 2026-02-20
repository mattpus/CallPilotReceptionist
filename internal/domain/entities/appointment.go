package entities

import (
	"time"

	"github.com/CallPilotReceptionist/internal/domain/errors"
)

type AppointmentStatus string

const (
	AppointmentStatusPending   AppointmentStatus = "pending"
	AppointmentStatusConfirmed AppointmentStatus = "confirmed"
	AppointmentStatusCancelled AppointmentStatus = "cancelled"
	AppointmentStatusCompleted AppointmentStatus = "completed"
)

type AppointmentRequest struct {
	ID              string            `json:"id"`
	CallID          string            `json:"call_id"`
	BusinessID      string            `json:"business_id"`
	CustomerName    string            `json:"customer_name"`
	CustomerPhone   string            `json:"customer_phone"`
	RequestedDate   *time.Time        `json:"requested_date,omitempty"`
	RequestedTime   string            `json:"requested_time,omitempty"`
	ServiceType     string            `json:"service_type,omitempty"`
	Notes           string            `json:"notes,omitempty"`
	Status          AppointmentStatus `json:"status"`
	ExtractedAt     time.Time         `json:"extracted_at"`
	ConfirmedAt     *time.Time        `json:"confirmed_at,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
}

func NewAppointmentRequest(
	callID, businessID, customerName, customerPhone string,
	requestedDate *time.Time, requestedTime, serviceType, notes string,
) (*AppointmentRequest, error) {
	if callID == "" {
		return nil, errors.NewValidationError("call_id is required")
	}
	if businessID == "" {
		return nil, errors.NewValidationError("business_id is required")
	}
	if customerPhone == "" {
		return nil, errors.NewValidationError("customer_phone is required")
	}

	now := time.Now()
	return &AppointmentRequest{
		CallID:        callID,
		BusinessID:    businessID,
		CustomerName:  customerName,
		CustomerPhone: customerPhone,
		RequestedDate: requestedDate,
		RequestedTime: requestedTime,
		ServiceType:   serviceType,
		Notes:         notes,
		Status:        AppointmentStatusPending,
		ExtractedAt:   now,
		CreatedAt:     now,
	}, nil
}

func (a *AppointmentRequest) Confirm() error {
	if a.Status != AppointmentStatusPending {
		return errors.NewValidationError("only pending appointments can be confirmed")
	}
	a.Status = AppointmentStatusConfirmed
	now := time.Now()
	a.ConfirmedAt = &now
	return nil
}

func (a *AppointmentRequest) Cancel() error {
	if a.Status == AppointmentStatusCompleted {
		return errors.NewValidationError("cannot cancel completed appointment")
	}
	a.Status = AppointmentStatusCancelled
	return nil
}

func (a *AppointmentRequest) Complete() error {
	if a.Status != AppointmentStatusConfirmed {
		return errors.NewValidationError("only confirmed appointments can be completed")
	}
	a.Status = AppointmentStatusCompleted
	return nil
}

func (a *AppointmentRequest) IsPending() bool {
	return a.Status == AppointmentStatusPending
}

func (a *AppointmentRequest) IsConfirmed() bool {
	return a.Status == AppointmentStatusConfirmed
}

func (a *AppointmentRequest) Validate() error {
	if a.CallID == "" {
		return errors.NewValidationError("call_id is required")
	}
	if a.BusinessID == "" {
		return errors.NewValidationError("business_id is required")
	}
	if a.CustomerPhone == "" {
		return errors.NewValidationError("customer_phone is required")
	}
	return nil
}
