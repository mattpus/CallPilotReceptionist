package services

import (
	"context"
	"time"

	"github.com/CallPilotReceptionist/internal/application/dto"
	"github.com/CallPilotReceptionist/internal/domain/entities"
	"github.com/CallPilotReceptionist/internal/domain/errors"
	"github.com/CallPilotReceptionist/internal/domain/providers"
	"github.com/CallPilotReceptionist/internal/infrastructure/database"
	"github.com/CallPilotReceptionist/pkg/logger"
)

type CallService struct {
	callRepo        database.CallRepository
	transcriptRepo  database.TranscriptRepository
	interactionRepo database.InteractionRepository
	voiceProvider   providers.VoiceProvider
	logger          *logger.Logger
}

func NewCallService(
	callRepo database.CallRepository,
	transcriptRepo database.TranscriptRepository,
	interactionRepo database.InteractionRepository,
	voiceProvider providers.VoiceProvider,
	log *logger.Logger,
) *CallService {
	return &CallService{
		callRepo:        callRepo,
		transcriptRepo:  transcriptRepo,
		interactionRepo: interactionRepo,
		voiceProvider:   voiceProvider,
		logger:          log,
	}
}

func (s *CallService) InitiateCall(ctx context.Context, businessID string, req dto.InitiateCallRequest) (*dto.CallResponse, error) {
	// Validate input
	if req.PhoneNumber == "" {
		return nil, errors.NewValidationError("phone_number is required")
	}

	// Create call entity
	call, err := entities.NewCall(businessID, req.PhoneNumber)
	if err != nil {
		return nil, err
	}

	// Save call to database
	if err := s.callRepo.Create(ctx, call); err != nil {
		s.logger.Error("Failed to create call record", err, map[string]interface{}{
			"business_id": businessID,
			"phone":       req.PhoneNumber,
		})
		return nil, err
	}

	// Initiate call with provider
	providerReq := providers.CallRequest{
		PhoneNumber: req.PhoneNumber,
		AssistantID: req.AssistantID,
		Metadata: map[string]interface{}{
			"call_id":     call.ID,
			"business_id": businessID,
		},
	}

	session, err := s.voiceProvider.InitiateCall(ctx, providerReq)
	if err != nil {
		s.logger.Error("Failed to initiate call with provider", err, map[string]interface{}{
			"call_id": call.ID,
		})
		// Update call status to failed
		call.UpdateStatus(entities.CallStatusFailed)
		s.callRepo.Update(ctx, call)
		return nil, err
	}

	// Update call with provider ID
	call.SetProviderCallID(session.ID)
	if err := s.callRepo.Update(ctx, call); err != nil {
		s.logger.Error("Failed to update call with provider ID", err, map[string]interface{}{
			"call_id":          call.ID,
			"provider_call_id": session.ID,
		})
	}

	s.logger.Info("Call initiated successfully", map[string]interface{}{
		"call_id":          call.ID,
		"provider_call_id": session.ID,
		"business_id":      businessID,
	})

	return s.mapCallToResponse(call), nil
}

func (s *CallService) HandleWebhook(ctx context.Context, payload []byte, signature string) error {
	// Process webhook from provider
	event, err := s.voiceProvider.HandleWebhook(ctx, payload, signature)
	if err != nil {
		s.logger.Error("Failed to parse webhook", err, nil)
		return err
	}

	s.logger.Info("Received webhook event", map[string]interface{}{
		"event_type": event.Type,
		"call_id":    event.CallID,
		"status":     event.Status,
	})

	// Get call from database using provider call ID
	call, err := s.callRepo.GetByProviderCallID(ctx, event.CallID)
	if err != nil {
		s.logger.Error("Call not found for webhook", err, map[string]interface{}{
			"provider_call_id": event.CallID,
		})
		return err
	}

	// Update call status based on event
	switch event.Type {
	case "call.started":
		call.UpdateStatus(entities.CallStatusInProgress)
	case "call.ended", "call.completed":
		call.UpdateStatus(entities.CallStatusCompleted)
		// Fetch and store transcript
		go s.fetchAndStoreTranscript(context.Background(), call.ID, event.CallID)
	case "call.failed":
		call.UpdateStatus(entities.CallStatusFailed)
	}

	// Update call in database
	if err := s.callRepo.Update(ctx, call); err != nil {
		s.logger.Error("Failed to update call from webhook", err, map[string]interface{}{
			"call_id": call.ID,
		})
		return err
	}

	return nil
}

func (s *CallService) GetCall(ctx context.Context, businessID, callID string) (*dto.CallResponse, error) {
	call, err := s.callRepo.GetByID(ctx, callID)
	if err != nil {
		return nil, err
	}

	// Verify business ownership
	if call.BusinessID != businessID {
		return nil, errors.NewForbiddenError("access denied to this call")
	}

	return s.mapCallToResponse(call), nil
}

func (s *CallService) ListCalls(ctx context.Context, businessID string, req dto.ListCallsRequest) (*dto.ListCallsResponse, error) {
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	calls, err := s.callRepo.GetByBusinessID(ctx, businessID, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	response := &dto.ListCallsResponse{
		Calls:  make([]dto.CallResponse, 0, len(calls)),
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	for _, call := range calls {
		response.Calls = append(response.Calls, *s.mapCallToResponse(call))
	}

	return response, nil
}

func (s *CallService) GetTranscript(ctx context.Context, businessID, callID string) (*dto.TranscriptResponse, error) {
	// Verify call belongs to business
	call, err := s.callRepo.GetByID(ctx, callID)
	if err != nil {
		return nil, err
	}

	if call.BusinessID != businessID {
		return nil, errors.NewForbiddenError("access denied to this call")
	}

	// Get transcript from database
	transcripts, err := s.transcriptRepo.GetByCallID(ctx, callID)
	if err != nil {
		return nil, err
	}

	response := &dto.TranscriptResponse{
		CallID:   callID,
		Messages: make([]dto.TranscriptMessageResponse, 0, len(transcripts)),
	}

	for _, t := range transcripts {
		response.Messages = append(response.Messages, dto.TranscriptMessageResponse{
			Role:      string(t.Role),
			Message:   t.Message,
			Timestamp: t.Timestamp.Format(time.RFC3339),
		})
	}

	return response, nil
}

func (s *CallService) fetchAndStoreTranscript(ctx context.Context, callID, providerCallID string) {
	transcript, err := s.voiceProvider.GetTranscript(ctx, providerCallID)
	if err != nil {
		s.logger.Error("Failed to fetch transcript from provider", err, map[string]interface{}{
			"call_id":          callID,
			"provider_call_id": providerCallID,
		})
		return
	}

	// Convert provider transcript to entities
	transcriptEntities := make([]*entities.Transcript, 0, len(transcript.Messages))
	for _, msg := range transcript.Messages {
		t := entities.NewTranscript(
			callID,
			entities.TranscriptRole(msg.Role),
			msg.Message,
			msg.Timestamp,
		)
		transcriptEntities = append(transcriptEntities, t)
	}

	// Store in database
	if err := s.transcriptRepo.CreateBatch(ctx, transcriptEntities); err != nil {
		s.logger.Error("Failed to store transcript", err, map[string]interface{}{
			"call_id": callID,
		})
	} else {
		s.logger.Info("Transcript stored successfully", map[string]interface{}{
			"call_id":      callID,
			"message_count": len(transcriptEntities),
		})
	}
}

func (s *CallService) mapCallToResponse(call *entities.Call) *dto.CallResponse {
	response := &dto.CallResponse{
		ID:             call.ID,
		BusinessID:     call.BusinessID,
		ProviderCallID: call.ProviderCallID,
		CallerPhone:    call.CallerPhone,
		Duration:       call.Duration,
		Status:         string(call.Status),
		Cost:           call.Cost,
		CreatedAt:      call.CreatedAt.Format(time.RFC3339),
	}

	if call.StartedAt != nil {
		startedAt := call.StartedAt.Format(time.RFC3339)
		response.StartedAt = &startedAt
	}

	if call.EndedAt != nil {
		endedAt := call.EndedAt.Format(time.RFC3339)
		response.EndedAt = &endedAt
	}

	return response
}
