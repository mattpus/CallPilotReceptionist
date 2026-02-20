package providers

import (
	"context"
	"time"
)

// CallRequest represents a request to initiate a call
type CallRequest struct {
	PhoneNumber      string                 `json:"phone_number"`
	AssistantID      string                 `json:"assistant_id,omitempty"`
	AssistantConfig  *AssistantConfig       `json:"assistant_config,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// CallSession represents an initiated call session
type CallSession struct {
	ID          string                 `json:"id"`
	Status      string                 `json:"status"`
	PhoneNumber string                 `json:"phone_number"`
	StartedAt   *time.Time             `json:"started_at,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// CallEvent represents a webhook event from the provider
type CallEvent struct {
	Type        string                 `json:"type"`
	CallID      string                 `json:"call_id"`
	Status      string                 `json:"status"`
	Timestamp   time.Time              `json:"timestamp"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

// CallDetails contains detailed information about a call
type CallDetails struct {
	ID           string                 `json:"id"`
	Status       string                 `json:"status"`
	PhoneNumber  string                 `json:"phone_number"`
	Duration     int                    `json:"duration"` // seconds
	Cost         float64                `json:"cost"`
	StartedAt    *time.Time             `json:"started_at,omitempty"`
	EndedAt      *time.Time             `json:"ended_at,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	ErrorMessage string                 `json:"error_message,omitempty"`
}

// Transcript contains the conversation transcript
type Transcript struct {
	CallID   string             `json:"call_id"`
	Messages []TranscriptMessage `json:"messages"`
}

// TranscriptMessage represents a single message in the transcript
type TranscriptMessage struct {
	Role      string    `json:"role"` // assistant, user, system
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// AssistantConfig defines the AI assistant configuration
type AssistantConfig struct {
	Name             string                 `json:"name"`
	Voice            string                 `json:"voice,omitempty"`
	Language         string                 `json:"language,omitempty"`
	Prompt           string                 `json:"prompt,omitempty"`
	FirstMessage     string                 `json:"first_message,omitempty"`
	Model            string                 `json:"model,omitempty"`
	Functions        []Function             `json:"functions,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// Function represents a callable function for the assistant
type Function struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// VoiceProvider defines the interface for voice AI providers
// This abstraction allows easy switching between providers (Vapi, Twilio, etc.)
type VoiceProvider interface {
	// InitiateCall starts a new outbound call
	InitiateCall(ctx context.Context, req CallRequest) (*CallSession, error)

	// HandleWebhook processes incoming webhook events
	HandleWebhook(ctx context.Context, payload []byte, signature string) (*CallEvent, error)

	// GetCallDetails retrieves detailed information about a specific call
	GetCallDetails(ctx context.Context, callID string) (*CallDetails, error)

	// GetTranscript fetches the conversation transcript for a call
	GetTranscript(ctx context.Context, callID string) (*Transcript, error)

	// UpdateAssistantConfig updates or creates an assistant configuration
	UpdateAssistantConfig(ctx context.Context, config AssistantConfig) (string, error)

	// GetAssistantConfig retrieves an assistant configuration by ID
	GetAssistantConfig(ctx context.Context, assistantID string) (*AssistantConfig, error)

	// DeleteAssistantConfig removes an assistant configuration
	DeleteAssistantConfig(ctx context.Context, assistantID string) error

	// ValidateWebhookSignature verifies the authenticity of webhook requests
	ValidateWebhookSignature(payload []byte, signature string) bool
}
