package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/CallPilotReceptionist/internal/application/dto"
	"github.com/CallPilotReceptionist/internal/domain/entities"
	"github.com/CallPilotReceptionist/internal/domain/providers"
	"github.com/CallPilotReceptionist/internal/infrastructure/database"
	"github.com/CallPilotReceptionist/pkg/logger"
)

// Extended mock for CallRepository with function fields
type testCallRepository struct {
	calls                   map[string]*entities.Call
	createFunc              func(ctx context.Context, call *entities.Call) error
	getByIDFunc             func(ctx context.Context, id string) (*entities.Call, error)
	getByProviderCallIDFunc func(ctx context.Context, providerCallID string) (*entities.Call, error)
	updateFunc              func(ctx context.Context, call *entities.Call) error
	listFunc                func(ctx context.Context, businessID string, limit, offset int) ([]*entities.Call, error)
}

func newTestCallRepository() *testCallRepository {
	return &testCallRepository{
		calls: make(map[string]*entities.Call),
	}
}

func (m *testCallRepository) Create(ctx context.Context, call *entities.Call) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, call)
	}
	if call.ID == "" {
		call.ID = "call-" + time.Now().Format("20060102150405")
	}
	m.calls[call.ID] = call
	return nil
}

func (m *testCallRepository) GetByID(ctx context.Context, id string) (*entities.Call, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	if call, ok := m.calls[id]; ok {
		return call, nil
	}
	return nil, errors.New("call not found")
}

func (m *testCallRepository) GetByProviderCallID(ctx context.Context, providerCallID string) (*entities.Call, error) {
	if m.getByProviderCallIDFunc != nil {
		return m.getByProviderCallIDFunc(ctx, providerCallID)
	}
	for _, call := range m.calls {
		if call.ProviderCallID == providerCallID {
			return call, nil
		}
	}
	return nil, errors.New("call not found")
}

func (m *testCallRepository) Update(ctx context.Context, call *entities.Call) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, call)
	}
	m.calls[call.ID] = call
	return nil
}

func (m *testCallRepository) ListByBusinessID(ctx context.Context, businessID string, limit, offset int) ([]*entities.Call, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, businessID, limit, offset)
	}
	var calls []*entities.Call
	for _, call := range m.calls {
		if call.BusinessID == businessID {
			calls = append(calls, call)
		}
	}
	return calls, nil
}

func (m *testCallRepository) GetByBusinessID(ctx context.Context, businessID string, limit, offset int) ([]*entities.Call, error) {
	return m.ListByBusinessID(ctx, businessID, limit, offset)
}

func (m *testCallRepository) Delete(ctx context.Context, id string) error {
	delete(m.calls, id)
	return nil
}

func (m *testCallRepository) GetByDateRange(ctx context.Context, businessID string, startDate, endDate time.Time) ([]*entities.Call, error) {
	return nil, errors.New("not implemented")
}

func (m *testCallRepository) GetStats(ctx context.Context, businessID string, startDate, endDate time.Time) (*database.CallStats, error) {
	return nil, errors.New("not implemented")
}

// Extended mock for TranscriptRepository
type testTranscriptRepository struct {
	transcripts     map[string][]*entities.Transcript
	createBatchFunc func(ctx context.Context, transcripts []*entities.Transcript) error
	getByCallIDFunc func(ctx context.Context, callID string) ([]*entities.Transcript, error)
}

func newTestTranscriptRepository() *testTranscriptRepository {
	return &testTranscriptRepository{
		transcripts: make(map[string][]*entities.Transcript),
	}
}

func (m *testTranscriptRepository) Create(ctx context.Context, transcript *entities.Transcript) error {
	if transcript.ID == "" {
		transcript.ID = "transcript-" + time.Now().Format("20060102150405")
	}
	m.transcripts[transcript.CallID] = append(m.transcripts[transcript.CallID], transcript)
	return nil
}

func (m *testTranscriptRepository) CreateBatch(ctx context.Context, transcripts []*entities.Transcript) error {
	if m.createBatchFunc != nil {
		return m.createBatchFunc(ctx, transcripts)
	}
	for _, t := range transcripts {
		if err := m.Create(ctx, t); err != nil {
			return err
		}
	}
	return nil
}

func (m *testTranscriptRepository) GetByCallID(ctx context.Context, callID string) ([]*entities.Transcript, error) {
	if m.getByCallIDFunc != nil {
		return m.getByCallIDFunc(ctx, callID)
	}
	if transcripts, ok := m.transcripts[callID]; ok {
		return transcripts, nil
	}
	return []*entities.Transcript{}, nil
}

func (m *testTranscriptRepository) Delete(ctx context.Context, id string) error {
	return errors.New("not implemented")
}

func (m *testTranscriptRepository) DeleteByCallID(ctx context.Context, callID string) error {
	return errors.New("not implemented")
}

func (m *testTranscriptRepository) GetByID(ctx context.Context, id string) (*entities.Transcript, error) {
	return nil, errors.New("not implemented")
}

// Extended mock for InteractionRepository
type testInteractionRepository struct {
	interactions map[string][]*entities.Interaction
}

func newTestInteractionRepository() *testInteractionRepository {
	return &testInteractionRepository{
		interactions: make(map[string][]*entities.Interaction),
	}
}

func (m *testInteractionRepository) Create(ctx context.Context, interaction *entities.Interaction) error {
	if interaction.ID == "" {
		interaction.ID = "interaction-" + time.Now().Format("20060102150405")
	}
	m.interactions[interaction.CallID] = append(m.interactions[interaction.CallID], interaction)
	return nil
}

func (m *testInteractionRepository) GetByCallID(ctx context.Context, callID string) ([]*entities.Interaction, error) {
	if interactions, ok := m.interactions[callID]; ok {
		return interactions, nil
	}
	return []*entities.Interaction{}, nil
}

func (m *testInteractionRepository) GetByID(ctx context.Context, id string) (*entities.Interaction, error) {
	return nil, errors.New("not implemented")
}

func (m *testInteractionRepository) Delete(ctx context.Context, id string) error {
	return errors.New("not implemented")
}

func (m *testInteractionRepository) List(ctx context.Context, businessID string, limit, offset int) ([]*entities.Interaction, error) {
	return nil, errors.New("not implemented")
}

// Extended mock for VoiceProvider
type testVoiceProvider struct {
	initiateCallFunc func(ctx context.Context, req providers.CallRequest) (*providers.CallSession, error)
	handleWebhookFunc func(ctx context.Context, payload []byte, signature string) (*providers.CallEvent, error)
	getTranscriptFunc func(ctx context.Context, callID string) (*providers.Transcript, error)
}

func (m *testVoiceProvider) InitiateCall(ctx context.Context, req providers.CallRequest) (*providers.CallSession, error) {
	if m.initiateCallFunc != nil {
		return m.initiateCallFunc(ctx, req)
	}
	return nil, errors.New("not implemented")
}

func (m *testVoiceProvider) HandleWebhook(ctx context.Context, payload []byte, signature string) (*providers.CallEvent, error) {
	if m.handleWebhookFunc != nil {
		return m.handleWebhookFunc(ctx, payload, signature)
	}
	return nil, errors.New("not implemented")
}

func (m *testVoiceProvider) GetCallDetails(ctx context.Context, callID string) (*providers.CallDetails, error) {
	return nil, errors.New("not implemented")
}

func (m *testVoiceProvider) GetTranscript(ctx context.Context, callID string) (*providers.Transcript, error) {
	if m.getTranscriptFunc != nil {
		return m.getTranscriptFunc(ctx, callID)
	}
	return nil, errors.New("not implemented")
}

func (m *testVoiceProvider) UpdateAssistantConfig(ctx context.Context, config providers.AssistantConfig) (string, error) {
	return "", errors.New("not implemented")
}

func (m *testVoiceProvider) GetAssistantConfig(ctx context.Context, assistantID string) (*providers.AssistantConfig, error) {
	return nil, errors.New("not implemented")
}

func (m *testVoiceProvider) DeleteAssistantConfig(ctx context.Context, assistantID string) error {
	return errors.New("not implemented")
}

func (m *testVoiceProvider) ValidateWebhookSignature(payload []byte, signature string) bool {
	return true
}

func TestCallService_InitiateCall(t *testing.T) {
	log := logger.New("info", "console")

	tests := []struct {
		name          string
		businessID    string
		request       dto.InitiateCallRequest
		setupMocks    func(*testCallRepository, *testVoiceProvider)
		expectedError bool
		validateResp  func(*testing.T, *dto.CallResponse)
	}{
		{
			name:       "successful call initiation",
			businessID: "business-123",
			request: dto.InitiateCallRequest{
				PhoneNumber: "+1234567890",
				AssistantID: "assistant-1",
				Metadata: map[string]interface{}{
					"customer": "John Doe",
				},
			},
			setupMocks: func(callRepo *testCallRepository, provider *testVoiceProvider) {
				provider.initiateCallFunc = func(ctx context.Context, req providers.CallRequest) (*providers.CallSession, error) {
					return &providers.CallSession{
						ID:          "provider-call-123",
						Status:      "initiated",
						PhoneNumber: "+1234567890",
					}, nil
				}
			},
			expectedError: false,
			validateResp: func(t *testing.T, resp *dto.CallResponse) {
				if resp.CallerPhone != "+1234567890" {
					t.Errorf("expected phone +1234567890, got %s", resp.CallerPhone)
				}
				if resp.Status == "" {
					t.Error("expected status to be set")
				}
				if resp.ProviderCallID != "provider-call-123" {
					t.Errorf("expected provider call ID provider-call-123, got %s", resp.ProviderCallID)
				}
			},
		},
		{
			name:       "missing phone number",
			businessID: "business-123",
			request: dto.InitiateCallRequest{
				PhoneNumber: "",
			},
			setupMocks:    func(callRepo *testCallRepository, provider *testVoiceProvider) {},
			expectedError: true,
		},
		{
			name:       "provider fails to initiate call",
			businessID: "business-123",
			request: dto.InitiateCallRequest{
				PhoneNumber: "+1234567890",
			},
			setupMocks: func(callRepo *testCallRepository, provider *testVoiceProvider) {
				provider.initiateCallFunc = func(ctx context.Context, req providers.CallRequest) (*providers.CallSession, error) {
					return nil, errors.New("provider API unavailable")
				}
			},
			expectedError: true,
		},
		{
			name:       "database error on create",
			businessID: "business-123",
			request: dto.InitiateCallRequest{
				PhoneNumber: "+1234567890",
			},
			setupMocks: func(callRepo *testCallRepository, provider *testVoiceProvider) {
				callRepo.createFunc = func(ctx context.Context, call *entities.Call) error {
					return errors.New("database connection failed")
				}
				provider.initiateCallFunc = func(ctx context.Context, req providers.CallRequest) (*providers.CallSession, error) {
					return &providers.CallSession{
						ID:          "provider-call-123",
						Status:      "initiated",
						PhoneNumber: "+1234567890",
					}, nil
				}
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			callRepo := newTestCallRepository()
			transcriptRepo := newTestTranscriptRepository()
			interactionRepo := newTestInteractionRepository()
			provider := &testVoiceProvider{}

			tt.setupMocks(callRepo, provider)

			service := NewCallService(
				callRepo,
				transcriptRepo,
				interactionRepo,
				provider,
				log,
			)

			response, err := service.InitiateCall(context.Background(), tt.businessID, tt.request)

			if tt.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if response == nil {
					t.Error("expected response but got nil")
					return
				}
				if tt.validateResp != nil {
					tt.validateResp(t, response)
				}
			}
		})
	}
}

func TestCallService_HandleWebhook(t *testing.T) {
	log := logger.New("info", "console")

	tests := []struct {
		name          string
		payload       []byte
		signature     string
		setupMocks    func(*testCallRepository, *testVoiceProvider)
		expectedError bool
	}{
		{
			name:      "successful webhook - call started",
			payload:   []byte(`{"type":"call.started","callId":"provider-123"}`),
			signature: "valid-signature",
			setupMocks: func(callRepo *testCallRepository, provider *testVoiceProvider) {
				// Create a call first
				call, _ := entities.NewCall("business-123", "+1234567890")
				call.ID = "call-123"
				call.ProviderCallID = "provider-123"
				call.Status = entities.CallStatusInitiated
				callRepo.calls[call.ID] = call

				provider.handleWebhookFunc = func(ctx context.Context, payload []byte, signature string) (*providers.CallEvent, error) {
					return &providers.CallEvent{
						Type:      "call.started",
						CallID:    "provider-123",
						Status:    "in_progress",
						Timestamp: time.Now(),
					}, nil
				}
			},
			expectedError: false,
		},
		{
			name:      "webhook - call ended",
			payload:   []byte(`{"type":"call.ended","callId":"provider-123"}`),
			signature: "valid-signature",
			setupMocks: func(callRepo *testCallRepository, provider *testVoiceProvider) {
				startTime := time.Now().Add(-2 * time.Minute)
				call, _ := entities.NewCall("business-123", "+1234567890")
				call.ID = "call-123"
				call.ProviderCallID = "provider-123"
				call.Status = entities.CallStatusInProgress
				call.StartedAt = &startTime
				callRepo.calls[call.ID] = call

				provider.handleWebhookFunc = func(ctx context.Context, payload []byte, signature string) (*providers.CallEvent, error) {
					return &providers.CallEvent{
						Type:      "call.ended",
						CallID:    "provider-123",
						Status:    "completed",
						Timestamp: time.Now(),
						Data: map[string]interface{}{
							"duration": 120,
						},
					}, nil
				}
			},
			expectedError: false,
		},
		{
			name:      "invalid webhook signature",
			payload:   []byte(`{"type":"call.started","callId":"provider-123"}`),
			signature: "invalid-signature",
			setupMocks: func(callRepo *testCallRepository, provider *testVoiceProvider) {
				provider.handleWebhookFunc = func(ctx context.Context, payload []byte, signature string) (*providers.CallEvent, error) {
					return nil, errors.New("invalid signature")
				}
			},
			expectedError: true,
		},
		{
			name:      "call not found in database",
			payload:   []byte(`{"type":"call.started","callId":"unknown-123"}`),
			signature: "valid-signature",
			setupMocks: func(callRepo *testCallRepository, provider *testVoiceProvider) {
				provider.handleWebhookFunc = func(ctx context.Context, payload []byte, signature string) (*providers.CallEvent, error) {
					return &providers.CallEvent{
						Type:      "call.started",
						CallID:    "unknown-123",
						Status:    "in_progress",
						Timestamp: time.Now(),
					}, nil
				}
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			callRepo := newTestCallRepository()
			transcriptRepo := newTestTranscriptRepository()
			interactionRepo := newTestInteractionRepository()
			provider := &testVoiceProvider{}

			tt.setupMocks(callRepo, provider)

			service := NewCallService(
				callRepo,
				transcriptRepo,
				interactionRepo,
				provider,
				log,
			)

			err := service.HandleWebhook(context.Background(), tt.payload, tt.signature)

			if tt.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestCallService_GetCall(t *testing.T) {
	log := logger.New("info", "console")

	businessID := "business-123"
	callID := "call-123"

	tests := []struct {
		name          string
		businessID    string
		callID        string
		setupMocks    func(*testCallRepository)
		expectedError bool
		validateResp  func(*testing.T, *dto.CallResponse)
	}{
		{
			name:       "successful get call",
			businessID: businessID,
			callID:     callID,
			setupMocks: func(callRepo *testCallRepository) {
				call, _ := entities.NewCall(businessID, "+1234567890")
				call.ID = callID
				call.ProviderCallID = "provider-123"
				call.Status = entities.CallStatusCompleted
				call.Duration = 120
				callRepo.calls[callID] = call
			},
			expectedError: false,
			validateResp: func(t *testing.T, resp *dto.CallResponse) {
				if resp.ID != callID {
					t.Errorf("expected call ID %s, got %s", callID, resp.ID)
				}
				if resp.Duration != 120 {
					t.Errorf("expected duration 120, got %d", resp.Duration)
				}
			},
		},
		{
			name:       "call not found",
			businessID: businessID,
			callID:     "non-existent",
			setupMocks: func(callRepo *testCallRepository) {
				// No calls in repository
			},
			expectedError: true,
		},
		{
			name:       "unauthorized - different business",
			businessID: "different-business",
			callID:     callID,
			setupMocks: func(callRepo *testCallRepository) {
				call, _ := entities.NewCall(businessID, "+1234567890")
				call.ID = callID
				callRepo.calls[callID] = call
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			callRepo := newTestCallRepository()
			tt.setupMocks(callRepo)

			service := NewCallService(
				callRepo,
				newTestTranscriptRepository(),
				newTestInteractionRepository(),
				&testVoiceProvider{},
				log,
			)

			response, err := service.GetCall(context.Background(), tt.businessID, tt.callID)

			if tt.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if response == nil {
					t.Error("expected response but got nil")
					return
				}
				if tt.validateResp != nil {
					tt.validateResp(t, response)
				}
			}
		})
	}
}

func TestCallService_GetTranscript(t *testing.T) {
log := logger.New("test", "debug")
callID := "call-123"
businessID := "business-123"

tests := []struct {
name          string
callID        string
setupMocks    func(*testCallRepository, *testTranscriptRepository)
expectedError bool
expectedCount int
}{
{
name:   "successful get transcript with messages",
callID: callID,
setupMocks: func(callRepo *testCallRepository, transcriptRepo *testTranscriptRepository) {
// Setup call
call, _ := entities.NewCall(businessID, "+1234567890")
call.ID = callID
call.ProviderCallID = "provider-123"
call.Status = "completed"
callRepo.calls[callID] = call

transcripts := []*entities.Transcript{
{
ID:        "t1",
CallID:    callID,
Role:      "assistant",
Message:   "Hello, how can I help you?",
Timestamp: time.Now(),
},
{
ID:        "t2",
CallID:    callID,
Role:      "user",
Message:   "I'd like to schedule an appointment",
Timestamp: time.Now().Add(1 * time.Second),
},
{
ID:        "t3",
CallID:    callID,
Role:      "assistant",
Message:   "I'd be happy to help with that.",
Timestamp: time.Now().Add(2 * time.Second),
},
}
transcriptRepo.transcripts[callID] = transcripts
},
expectedError: false,
expectedCount: 3,
},
{
name:   "no transcript found",
callID: callID,
setupMocks: func(callRepo *testCallRepository, transcriptRepo *testTranscriptRepository) {
// Setup call
call, _ := entities.NewCall(businessID, "+1234567890")
call.ID = callID
call.ProviderCallID = "provider-123"
call.Status = "completed"
callRepo.calls[callID] = call
// No transcripts
},
expectedError: false,
expectedCount: 0,
},
{
name:   "database error",
callID: callID,
setupMocks: func(callRepo *testCallRepository, transcriptRepo *testTranscriptRepository) {
// Setup call
call, _ := entities.NewCall(businessID, "+1234567890")
call.ID = callID
call.ProviderCallID = "provider-123"
call.Status = "completed"
callRepo.calls[callID] = call

transcriptRepo.getByCallIDFunc = func(ctx context.Context, cid string) ([]*entities.Transcript, error) {
return nil, errors.New("database connection lost")
}
},
expectedError: true,
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
callRepo := newTestCallRepository()
transcriptRepo := newTestTranscriptRepository()
tt.setupMocks(callRepo, transcriptRepo)

service := NewCallService(
callRepo,
transcriptRepo,
newTestInteractionRepository(),
&testVoiceProvider{},
log,
)

response, err := service.GetTranscript(context.Background(), businessID, tt.callID)

if tt.expectedError {
if err == nil {
t.Error("expected error but got none")
}
} else {
if err != nil {
t.Errorf("unexpected error: %v", err)
return
}
if response == nil {
t.Error("expected response but got nil")
return
}
if len(response.Messages) != tt.expectedCount {
t.Errorf("expected %d messages, got %d", tt.expectedCount, len(response.Messages))
}
}
})
}
}

