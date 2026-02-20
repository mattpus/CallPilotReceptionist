package entities

import (
	"time"

	"github.com/CallPilotReceptionist/internal/domain/errors"
)

type InteractionType string

const (
	InteractionTypeAppointmentRequest InteractionType = "appointment_request"
	InteractionTypeQuestion           InteractionType = "question"
	InteractionTypeComplaint          InteractionType = "complaint"
	InteractionTypeInformation        InteractionType = "information"
	InteractionTypeGreeting           InteractionType = "greeting"
	InteractionTypeFarewell           InteractionType = "farewell"
	InteractionTypeOther              InteractionType = "other"
)

type Interaction struct {
	ID        string                 `json:"id"`
	CallID    string                 `json:"call_id"`
	Type      InteractionType        `json:"type"`
	Content   map[string]interface{} `json:"content"`
	Timestamp time.Time              `json:"timestamp"`
	CreatedAt time.Time              `json:"created_at"`
}

func NewInteraction(callID string, interactionType InteractionType, content map[string]interface{}) (*Interaction, error) {
	if callID == "" {
		return nil, errors.NewValidationError("call_id is required")
	}

	validTypes := map[InteractionType]bool{
		InteractionTypeAppointmentRequest: true,
		InteractionTypeQuestion:           true,
		InteractionTypeComplaint:          true,
		InteractionTypeInformation:        true,
		InteractionTypeGreeting:           true,
		InteractionTypeFarewell:           true,
		InteractionTypeOther:              true,
	}

	if !validTypes[interactionType] {
		return nil, errors.NewValidationError("invalid interaction type")
	}

	now := time.Now()
	return &Interaction{
		CallID:    callID,
		Type:      interactionType,
		Content:   content,
		Timestamp: now,
		CreatedAt: now,
	}, nil
}

func (i *Interaction) IsAppointmentRequest() bool {
	return i.Type == InteractionTypeAppointmentRequest
}

func (i *Interaction) Validate() error {
	if i.CallID == "" {
		return errors.NewValidationError("call_id is required")
	}
	if i.Type == "" {
		return errors.NewValidationError("interaction type is required")
	}
	return nil
}
