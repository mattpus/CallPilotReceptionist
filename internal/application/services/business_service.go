package services

import (
	"context"
	"errors"
	"time"

	"github.com/CallPilotReceptionist/internal/application/dto"
	"github.com/CallPilotReceptionist/internal/infrastructure/database"
	"github.com/CallPilotReceptionist/pkg/logger"
)

var ErrBusinessNotFound = errors.New("business not found")

type BusinessService struct {
	businessRepo database.BusinessRepository
	logger       *logger.Logger
}

func NewBusinessService(
	businessRepo database.BusinessRepository,
	log *logger.Logger,
) *BusinessService {
	return &BusinessService{
		businessRepo: businessRepo,
		logger:       log,
	}
}

func (s *BusinessService) GetBusiness(ctx context.Context, businessID string) (*dto.BusinessResponse, error) {
	business, err := s.businessRepo.GetByID(ctx, businessID)
	if err != nil {
		return nil, err
	}

	if business == nil {
		return nil, ErrBusinessNotFound
	}

	return &dto.BusinessResponse{
		ID:        business.ID,
		Name:      business.Name,
		Type:      business.Type,
		Phone:     business.Phone,
		Settings:  business.Settings,
		CreatedAt: business.CreatedAt.Format(time.RFC3339),
		UpdatedAt: business.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *BusinessService) UpdateBusiness(ctx context.Context, businessID string, req dto.UpdateBusinessRequest) (*dto.BusinessResponse, error) {
	// Get existing business
	business, err := s.businessRepo.GetByID(ctx, businessID)
	if err != nil {
		return nil, err
	}

	if business == nil {
		return nil, ErrBusinessNotFound
	}

	// Update fields
	if err := business.Update(req.Name, req.Type, req.Phone, req.Settings); err != nil {
		return nil, err
	}

	// Save to database
	if err := s.businessRepo.Update(ctx, business); err != nil {
		s.logger.Error("Failed to update business", err, map[string]interface{}{
			"business_id": businessID,
		})
		return nil, err
	}

	s.logger.Info("Business updated successfully", map[string]interface{}{
		"business_id": businessID,
	})

	return &dto.BusinessResponse{
		ID:        business.ID,
		Name:      business.Name,
		Type:      business.Type,
		Phone:     business.Phone,
		Settings:  business.Settings,
		CreatedAt: business.CreatedAt.Format(time.RFC3339),
		UpdatedAt: business.UpdatedAt.Format(time.RFC3339),
	}, nil
}
