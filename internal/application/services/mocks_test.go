package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/CallPilotReceptionist/internal/domain/entities"
	"github.com/CallPilotReceptionist/internal/domain/providers"
)

// ========== Mock VoiceProvider ==========

type mockVoiceProvider struct {
	initiateCallFunc              func(ctx context.Context, request providers.CallRequest) (*providers.CallSession, error)
	handleWebhookFunc             func(ctx context.Context, payload []byte, signature string) (*providers.CallEvent, error)
	getCallDetailsFunc            func(ctx context.Context, callID string) (*providers.CallDetails, error)
	getTranscriptFunc             func(ctx context.Context, callID string) (*providers.Transcript, error)
	updateAssistantConfigFunc     func(ctx context.Context, config providers.AssistantConfig) (string, error)
	getAssistantConfigFunc        func(ctx context.Context, assistantID string) (*providers.AssistantConfig, error)
	deleteAssistantConfigFunc     func(ctx context.Context, assistantID string) error
	validateWebhookSignatureFunc  func(payload []byte, signature string) bool
}

func (m *mockVoiceProvider) InitiateCall(ctx context.Context, request providers.CallRequest) (*providers.CallSession, error) {
	if m.initiateCallFunc != nil {
		return m.initiateCallFunc(ctx, request)
	}
	return nil, errors.New("not implemented")
}

func (m *mockVoiceProvider) HandleWebhook(ctx context.Context, payload []byte, signature string) (*providers.CallEvent, error) {
	if m.handleWebhookFunc != nil {
		return m.handleWebhookFunc(ctx, payload, signature)
	}
	return nil, errors.New("not implemented")
}

func (m *mockVoiceProvider) GetCallDetails(ctx context.Context, callID string) (*providers.CallDetails, error) {
	if m.getCallDetailsFunc != nil {
		return m.getCallDetailsFunc(ctx, callID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockVoiceProvider) GetTranscript(ctx context.Context, callID string) (*providers.Transcript, error) {
	if m.getTranscriptFunc != nil {
		return m.getTranscriptFunc(ctx, callID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockVoiceProvider) UpdateAssistantConfig(ctx context.Context, config providers.AssistantConfig) (string, error) {
	if m.updateAssistantConfigFunc != nil {
		return m.updateAssistantConfigFunc(ctx, config)
	}
	return "", errors.New("not implemented")
}

func (m *mockVoiceProvider) GetAssistantConfig(ctx context.Context, assistantID string) (*providers.AssistantConfig, error) {
	if m.getAssistantConfigFunc != nil {
		return m.getAssistantConfigFunc(ctx, assistantID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockVoiceProvider) DeleteAssistantConfig(ctx context.Context, assistantID string) error {
	if m.deleteAssistantConfigFunc != nil {
		return m.deleteAssistantConfigFunc(ctx, assistantID)
	}
	return errors.New("not implemented")
}

func (m *mockVoiceProvider) ValidateWebhookSignature(payload []byte, signature string) bool {
	if m.validateWebhookSignatureFunc != nil {
		return m.validateWebhookSignatureFunc(payload, signature)
	}
	return false
}

// ========== Mock CallRepository ==========

type mockCallRepository struct {
	createFunc               func(ctx context.Context, call *entities.Call) error
	getByIDFunc              func(ctx context.Context, id uuid.UUID) (*entities.Call, error)
	getByProviderCallIDFunc  func(ctx context.Context, providerCallID string) (*entities.Call, error)
	updateFunc               func(ctx context.Context, call *entities.Call) error
	listByBusinessIDFunc     func(ctx context.Context, businessID uuid.UUID, limit, offset int) ([]*entities.Call, error)
	deleteFunc               func(ctx context.Context, id uuid.UUID) error
}

func (m *mockCallRepository) Create(ctx context.Context, call *entities.Call) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, call)
	}
	return nil
}

func (m *mockCallRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Call, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not found")
}

func (m *mockCallRepository) GetByProviderCallID(ctx context.Context, providerCallID string) (*entities.Call, error) {
	if m.getByProviderCallIDFunc != nil {
		return m.getByProviderCallIDFunc(ctx, providerCallID)
	}
	return nil, errors.New("not found")
}

func (m *mockCallRepository) Update(ctx context.Context, call *entities.Call) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, call)
	}
	return nil
}

func (m *mockCallRepository) ListByBusinessID(ctx context.Context, businessID uuid.UUID, limit, offset int) ([]*entities.Call, error) {
	if m.listByBusinessIDFunc != nil {
		return m.listByBusinessIDFunc(ctx, businessID, limit, offset)
	}
	return []*entities.Call{}, nil
}

func (m *mockCallRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

// ========== Mock TranscriptRepository ==========

type mockTranscriptRepository struct {
	createFunc      func(ctx context.Context, transcript *entities.Transcript) error
	createBatchFunc func(ctx context.Context, transcripts []*entities.Transcript) error
	getByCallIDFunc func(ctx context.Context, callID uuid.UUID) ([]*entities.Transcript, error)
}

func (m *mockTranscriptRepository) Create(ctx context.Context, transcript *entities.Transcript) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, transcript)
	}
	return nil
}

func (m *mockTranscriptRepository) CreateBatch(ctx context.Context, transcripts []*entities.Transcript) error {
	if m.createBatchFunc != nil {
		return m.createBatchFunc(ctx, transcripts)
	}
	return nil
}

func (m *mockTranscriptRepository) GetByCallID(ctx context.Context, callID uuid.UUID) ([]*entities.Transcript, error) {
	if m.getByCallIDFunc != nil {
		return m.getByCallIDFunc(ctx, callID)
	}
	return []*entities.Transcript{}, nil
}

// ========== Mock InteractionRepository ==========

type mockInteractionRepository struct {
	createFunc       func(ctx context.Context, interaction *entities.Interaction) error
	getByCallIDFunc  func(ctx context.Context, callID uuid.UUID) ([]*entities.Interaction, error)
}

func (m *mockInteractionRepository) Create(ctx context.Context, interaction *entities.Interaction) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, interaction)
	}
	return nil
}

func (m *mockInteractionRepository) GetByCallID(ctx context.Context, callID uuid.UUID) ([]*entities.Interaction, error) {
	if m.getByCallIDFunc != nil {
		return m.getByCallIDFunc(ctx, callID)
	}
	return []*entities.Interaction{}, nil
}

// ========== Mock AppointmentRepository ==========

type mockAppointmentRepository struct {
	createFunc           func(ctx context.Context, appointment *entities.AppointmentRequest) error
	getByIDFunc          func(ctx context.Context, id uuid.UUID) (*entities.AppointmentRequest, error)
	updateFunc           func(ctx context.Context, appointment *entities.AppointmentRequest) error
	listByBusinessIDFunc func(ctx context.Context, businessID uuid.UUID, filters map[string]interface{}, limit, offset int) ([]*entities.AppointmentRequest, error)
}

func (m *mockAppointmentRepository) Create(ctx context.Context, appointment *entities.AppointmentRequest) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, appointment)
	}
	return nil
}

func (m *mockAppointmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.AppointmentRequest, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not found")
}

func (m *mockAppointmentRepository) Update(ctx context.Context, appointment *entities.AppointmentRequest) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, appointment)
	}
	return nil
}

func (m *mockAppointmentRepository) ListByBusinessID(ctx context.Context, businessID uuid.UUID, filters map[string]interface{}, limit, offset int) ([]*entities.AppointmentRequest, error) {
	if m.listByBusinessIDFunc != nil {
		return m.listByBusinessIDFunc(ctx, businessID, filters, limit, offset)
	}
	return []*entities.AppointmentRequest{}, nil
}
