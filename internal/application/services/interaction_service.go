package services

import (
	"context"
	"time"

	"github.com/CallPilotReceptionist/internal/application/dto"
	"github.com/CallPilotReceptionist/internal/domain/entities"
	"github.com/CallPilotReceptionist/internal/domain/errors"
	"github.com/CallPilotReceptionist/internal/infrastructure/database"
	"github.com/CallPilotReceptionist/pkg/logger"
)

type InteractionService struct {
	interactionRepo database.InteractionRepository
	appointmentRepo database.AppointmentRepository
	callRepo        database.CallRepository
	logger          *logger.Logger
}

func NewInteractionService(
	interactionRepo database.InteractionRepository,
	appointmentRepo database.AppointmentRepository,
	callRepo database.CallRepository,
	log *logger.Logger,
) *InteractionService {
	return &InteractionService{
		interactionRepo: interactionRepo,
		appointmentRepo: appointmentRepo,
		callRepo:        callRepo,
		logger:          log,
	}
}

func (s *InteractionService) GetCallInteractions(ctx context.Context, businessID, callID string) ([]dto.InteractionResponse, error) {
	// Verify call belongs to business
	call, err := s.callRepo.GetByID(ctx, callID)
	if err != nil {
		return nil, err
	}

	if call.BusinessID != businessID {
		return nil, errors.NewForbiddenError("access denied to this call")
	}

	// Get interactions
	interactions, err := s.interactionRepo.GetByCallID(ctx, callID)
	if err != nil {
		return nil, err
	}

	response := make([]dto.InteractionResponse, 0, len(interactions))
	for _, interaction := range interactions {
		response = append(response, dto.InteractionResponse{
			ID:        interaction.ID,
			CallID:    interaction.CallID,
			Type:      string(interaction.Type),
			Content:   interaction.Content,
			Timestamp: interaction.Timestamp.Format(time.RFC3339),
			CreatedAt: interaction.CreatedAt.Format(time.RFC3339),
		})
	}

	return response, nil
}

func (s *InteractionService) ListInteractions(ctx context.Context, businessID string, limit, offset int) (*dto.ListInteractionsResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	interactions, err := s.interactionRepo.List(ctx, businessID, limit, offset)
	if err != nil {
		return nil, err
	}

	response := &dto.ListInteractionsResponse{
		Interactions: make([]dto.InteractionResponse, 0, len(interactions)),
		Limit:        limit,
		Offset:       offset,
	}

	for _, interaction := range interactions {
		response.Interactions = append(response.Interactions, dto.InteractionResponse{
			ID:        interaction.ID,
			CallID:    interaction.CallID,
			Type:      string(interaction.Type),
			Content:   interaction.Content,
			Timestamp: interaction.Timestamp.Format(time.RFC3339),
			CreatedAt: interaction.CreatedAt.Format(time.RFC3339),
		})
	}

	return response, nil
}

func (s *InteractionService) GetAppointments(ctx context.Context, businessID string, limit, offset int) ([]dto.AppointmentResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	appointments, err := s.appointmentRepo.GetByBusinessID(ctx, businessID, limit, offset)
	if err != nil {
		return nil, err
	}

	response := make([]dto.AppointmentResponse, 0, len(appointments))
	for _, apt := range appointments {
		aptResponse := dto.AppointmentResponse{
			ID:            apt.ID,
			CallID:        apt.CallID,
			BusinessID:    apt.BusinessID,
			CustomerName:  apt.CustomerName,
			CustomerPhone: apt.CustomerPhone,
			RequestedTime: apt.RequestedTime,
			ServiceType:   apt.ServiceType,
			Notes:         apt.Notes,
			Status:        string(apt.Status),
			ExtractedAt:   apt.ExtractedAt.Format(time.RFC3339),
			CreatedAt:     apt.CreatedAt.Format(time.RFC3339),
		}

		if apt.RequestedDate != nil {
			dateStr := apt.RequestedDate.Format("2006-01-02")
			aptResponse.RequestedDate = &dateStr
		}

		if apt.ConfirmedAt != nil {
			confirmedStr := apt.ConfirmedAt.Format(time.RFC3339)
			aptResponse.ConfirmedAt = &confirmedStr
		}

		response = append(response, aptResponse)
	}

	return response, nil
}

func (s *InteractionService) UpdateAppointmentStatus(ctx context.Context, businessID, appointmentID string, req dto.UpdateAppointmentRequest) (*dto.AppointmentResponse, error) {
	// Get appointment
	apt, err := s.appointmentRepo.GetByID(ctx, appointmentID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if apt.BusinessID != businessID {
		return nil, errors.NewForbiddenError("access denied to this appointment")
	}

	// Update status
	switch entities.AppointmentStatus(req.Status) {
	case entities.AppointmentStatusConfirmed:
		if err := apt.Confirm(); err != nil {
			return nil, err
		}
	case entities.AppointmentStatusCancelled:
		if err := apt.Cancel(); err != nil {
			return nil, err
		}
	case entities.AppointmentStatusCompleted:
		if err := apt.Complete(); err != nil {
			return nil, err
		}
	default:
		return nil, errors.NewValidationError("invalid appointment status")
	}

	// Save to database
	if err := s.appointmentRepo.Update(ctx, apt); err != nil {
		s.logger.Error("Failed to update appointment", err, map[string]interface{}{
			"appointment_id": appointmentID,
		})
		return nil, err
	}

	s.logger.Info("Appointment status updated", map[string]interface{}{
		"appointment_id": appointmentID,
		"new_status":     req.Status,
	})

	// Build response
	aptResponse := dto.AppointmentResponse{
		ID:            apt.ID,
		CallID:        apt.CallID,
		BusinessID:    apt.BusinessID,
		CustomerName:  apt.CustomerName,
		CustomerPhone: apt.CustomerPhone,
		RequestedTime: apt.RequestedTime,
		ServiceType:   apt.ServiceType,
		Notes:         apt.Notes,
		Status:        string(apt.Status),
		ExtractedAt:   apt.ExtractedAt.Format(time.RFC3339),
		CreatedAt:     apt.CreatedAt.Format(time.RFC3339),
	}

	if apt.RequestedDate != nil {
		dateStr := apt.RequestedDate.Format("2006-01-02")
		aptResponse.RequestedDate = &dateStr
	}

	if apt.ConfirmedAt != nil {
		confirmedStr := apt.ConfirmedAt.Format(time.RFC3339)
		aptResponse.ConfirmedAt = &confirmedStr
	}

	return &aptResponse, nil
}
