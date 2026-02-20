package services

import (
	"context"
	"time"

	"github.com/CallPilotReceptionist/internal/application/dto"
	"github.com/CallPilotReceptionist/internal/infrastructure/database"
	"github.com/CallPilotReceptionist/pkg/logger"
)

type AnalyticsService struct {
	callRepo        database.CallRepository
	appointmentRepo database.AppointmentRepository
	logger          *logger.Logger
}

func NewAnalyticsService(
	callRepo database.CallRepository,
	appointmentRepo database.AppointmentRepository,
	log *logger.Logger,
) *AnalyticsService {
	return &AnalyticsService{
		callRepo:        callRepo,
		appointmentRepo: appointmentRepo,
		logger:          log,
	}
}

func (s *AnalyticsService) GetOverview(ctx context.Context, businessID string, days int) (*dto.AnalyticsOverviewResponse, error) {
	if days <= 0 {
		days = 30
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	// Get call statistics
	stats, err := s.callRepo.GetStats(ctx, businessID, startDate, endDate)
	if err != nil {
		s.logger.Error("Failed to get call stats", err, map[string]interface{}{
			"business_id": businessID,
		})
		return nil, err
	}

	// Get pending appointments count
	pendingAppointments, err := s.appointmentRepo.GetPendingAppointments(ctx, businessID)
	if err != nil {
		s.logger.Error("Failed to get pending appointments", err, map[string]interface{}{
			"business_id": businessID,
		})
		return nil, err
	}

	return &dto.AnalyticsOverviewResponse{
		TotalCalls:          stats.TotalCalls,
		CompletedCalls:      stats.CompletedCalls,
		FailedCalls:         stats.FailedCalls,
		TotalDuration:       stats.TotalDuration,
		AverageDuration:     stats.AverageDuration,
		TotalCost:           stats.TotalCost,
		PendingAppointments: len(pendingAppointments),
	}, nil
}

func (s *AnalyticsService) GetCallVolume(ctx context.Context, businessID string, days int) (*dto.CallVolumeResponse, error) {
	if days <= 0 {
		days = 30
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	// Get all calls in date range
	calls, err := s.callRepo.GetByDateRange(ctx, businessID, startDate, endDate)
	if err != nil {
		s.logger.Error("Failed to get calls by date range", err, map[string]interface{}{
			"business_id": businessID,
		})
		return nil, err
	}

	// Group calls by date
	volumeByDate := make(map[string]int)
	for _, call := range calls {
		dateKey := call.CreatedAt.Format("2006-01-02")
		volumeByDate[dateKey]++
	}

	// Convert to response format
	response := &dto.CallVolumeResponse{
		Data: make([]dto.CallVolumeData, 0),
	}

	// Fill in all dates even if no calls
	for d := startDate; d.Before(endDate) || d.Equal(endDate); d = d.AddDate(0, 0, 1) {
		dateKey := d.Format("2006-01-02")
		count := volumeByDate[dateKey]
		response.Data = append(response.Data, dto.CallVolumeData{
			Date:  dateKey,
			Count: count,
		})
	}

	return response, nil
}
