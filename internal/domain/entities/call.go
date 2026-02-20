package entities

import (
	"time"

	"github.com/CallPilotReceptionist/internal/domain/errors"
)

type CallStatus string

const (
	CallStatusInitiated  CallStatus = "initiated"
	CallStatusRinging    CallStatus = "ringing"
	CallStatusInProgress CallStatus = "in_progress"
	CallStatusCompleted  CallStatus = "completed"
	CallStatusFailed     CallStatus = "failed"
	CallStatusNoAnswer   CallStatus = "no_answer"
	CallStatusBusy       CallStatus = "busy"
)

type Call struct {
	ID             string     `json:"id"`
	BusinessID     string     `json:"business_id"`
	ProviderCallID string     `json:"provider_call_id"`
	CallerPhone    string     `json:"caller_phone"`
	Duration       int        `json:"duration"` // in seconds
	Status         CallStatus `json:"status"`
	Cost           float64    `json:"cost"`
	StartedAt      *time.Time `json:"started_at,omitempty"`
	EndedAt        *time.Time `json:"ended_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

func NewCall(businessID, callerPhone string) (*Call, error) {
	if businessID == "" {
		return nil, errors.NewValidationError("business_id is required")
	}
	if callerPhone == "" {
		return nil, errors.NewValidationError("caller_phone is required")
	}

	return &Call{
		BusinessID:  businessID,
		CallerPhone: callerPhone,
		Status:      CallStatusInitiated,
		CreatedAt:   time.Now(),
	}, nil
}

func (c *Call) UpdateStatus(status CallStatus) error {
	validStatuses := map[CallStatus]bool{
		CallStatusInitiated:  true,
		CallStatusRinging:    true,
		CallStatusInProgress: true,
		CallStatusCompleted:  true,
		CallStatusFailed:     true,
		CallStatusNoAnswer:   true,
		CallStatusBusy:       true,
	}

	if !validStatuses[status] {
		return errors.NewValidationError("invalid call status")
	}

	c.Status = status

	if status == CallStatusInProgress && c.StartedAt == nil {
		now := time.Now()
		c.StartedAt = &now
	}

	if (status == CallStatusCompleted || status == CallStatusFailed || 
		status == CallStatusNoAnswer || status == CallStatusBusy) && c.EndedAt == nil {
		now := time.Now()
		c.EndedAt = &now
		
		if c.StartedAt != nil {
			c.Duration = int(c.EndedAt.Sub(*c.StartedAt).Seconds())
		}
	}

	return nil
}

func (c *Call) SetProviderCallID(providerCallID string) {
	c.ProviderCallID = providerCallID
}

func (c *Call) SetCost(cost float64) {
	c.Cost = cost
}

func (c *Call) IsCompleted() bool {
	return c.Status == CallStatusCompleted || c.Status == CallStatusFailed || 
		c.Status == CallStatusNoAnswer || c.Status == CallStatusBusy
}

func (c *Call) Validate() error {
	if c.BusinessID == "" {
		return errors.NewValidationError("business_id is required")
	}
	if c.CallerPhone == "" {
		return errors.NewValidationError("caller_phone is required")
	}
	return nil
}
