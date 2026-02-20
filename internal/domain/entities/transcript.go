package entities

import (
	"time"
)

type TranscriptRole string

const (
	TranscriptRoleAssistant TranscriptRole = "assistant"
	TranscriptRoleUser      TranscriptRole = "user"
	TranscriptRoleSystem    TranscriptRole = "system"
)

type Transcript struct {
	ID        string         `json:"id"`
	CallID    string         `json:"call_id"`
	Role      TranscriptRole `json:"role"`
	Message   string         `json:"message"`
	Timestamp time.Time      `json:"timestamp"`
	CreatedAt time.Time      `json:"created_at"`
}

func NewTranscript(callID string, role TranscriptRole, message string, timestamp time.Time) *Transcript {
	return &Transcript{
		CallID:    callID,
		Role:      role,
		Message:   message,
		Timestamp: timestamp,
		CreatedAt: time.Now(),
	}
}

func (t *Transcript) IsFromAssistant() bool {
	return t.Role == TranscriptRoleAssistant
}

func (t *Transcript) IsFromUser() bool {
	return t.Role == TranscriptRoleUser
}
