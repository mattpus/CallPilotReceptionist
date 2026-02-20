package services

import (
	"context"
	"testing"

	"github.com/CallPilotReceptionist/internal/application/dto"
	"github.com/CallPilotReceptionist/internal/domain/entities"
	"github.com/CallPilotReceptionist/pkg/logger"
)

func TestBusinessService_GetBusiness(t *testing.T) {
	log := logger.New("info", "console")

	businessID := "business-123"

	tests := []struct {
		name          string
		businessID    string
		setupMocks    func(*mockBusinessRepository)
		expectedError bool
		validateResp  func(*testing.T, *dto.BusinessResponse)
	}{
		{
			name:       "successful get business",
			businessID: businessID,
			setupMocks: func(repo *mockBusinessRepository) {
				business := &entities.Business{
					ID:    businessID,
					Name:  "Smith Dental",
					Type:  "dentist",
					Phone: "+1234567890",
					Settings: map[string]interface{}{
						"timezone": "America/New_York",
					},
				}
				repo.businesses[businessID] = business
			},
			expectedError: false,
			validateResp: func(t *testing.T, resp *dto.BusinessResponse) {
				if resp.Name != "Smith Dental" {
					t.Errorf("expected name 'Smith Dental', got %s", resp.Name)
				}
				if resp.Type != "dentist" {
					t.Errorf("expected type 'dentist', got %s", resp.Type)
				}
				if resp.Phone != "+1234567890" {
					t.Errorf("expected phone '+1234567890', got %s", resp.Phone)
				}
			},
		},
		{
			name:       "business not found",
			businessID: "non-existent",
			setupMocks: func(repo *mockBusinessRepository) {
				// No business in repo
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockBusinessRepository{
				businesses: make(map[string]*entities.Business),
			}
			tt.setupMocks(repo)

			service := NewBusinessService(repo, log)

			response, err := service.GetBusiness(context.Background(), tt.businessID)

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

func TestBusinessService_UpdateBusiness(t *testing.T) {
	log := logger.New("info", "console")

	businessID := "business-123"

	tests := []struct {
		name          string
		businessID    string
		request       dto.UpdateBusinessRequest
		setupMocks    func(*mockBusinessRepository)
		expectedError bool
		validateResp  func(*testing.T, *dto.BusinessResponse)
	}{
		{
			name:       "successful update - name only",
			businessID: businessID,
			request: dto.UpdateBusinessRequest{
				Name: "Updated Dental Clinic",
			},
			setupMocks: func(repo *mockBusinessRepository) {
				business := &entities.Business{
					ID:    businessID,
					Name:  "Smith Dental",
					Type:  "dentist",
					Phone: "+1234567890",
				}
				repo.businesses[businessID] = business
			},
			expectedError: false,
			validateResp: func(t *testing.T, resp *dto.BusinessResponse) {
				if resp.Name != "Updated Dental Clinic" {
					t.Errorf("expected name 'Updated Dental Clinic', got %s", resp.Name)
				}
			},
		},
		{
			name:       "successful update - all fields",
			businessID: businessID,
			request: dto.UpdateBusinessRequest{
				Name:  "New Name",
				Type:  "clinic",
				Phone: "+9876543210",
				Settings: map[string]interface{}{
					"timezone":       "America/Los_Angeles",
					"workingHours":   "9-5",
					"appointmentGap": 30,
				},
			},
			setupMocks: func(repo *mockBusinessRepository) {
				business := &entities.Business{
					ID:    businessID,
					Name:  "Smith Dental",
					Type:  "dentist",
					Phone: "+1234567890",
				}
				repo.businesses[businessID] = business
			},
			expectedError: false,
			validateResp: func(t *testing.T, resp *dto.BusinessResponse) {
				if resp.Name != "New Name" {
					t.Errorf("expected name 'New Name', got %s", resp.Name)
				}
				if resp.Type != "clinic" {
					t.Errorf("expected type 'clinic', got %s", resp.Type)
				}
				if resp.Phone != "+9876543210" {
					t.Errorf("expected phone '+9876543210', got %s", resp.Phone)
				}
				if resp.Settings == nil {
					t.Error("expected settings to be set")
				}
			},
		},
		{
			name:       "business not found",
			businessID: "non-existent",
			request: dto.UpdateBusinessRequest{
				Name: "Updated Name",
			},
			setupMocks: func(repo *mockBusinessRepository) {
				// No business in repo
			},
			expectedError: true,
		},
		{
			name:       "update with empty name (no change)",
			businessID: businessID,
			request: dto.UpdateBusinessRequest{
				Phone: "+9876543210",
			},
			setupMocks: func(repo *mockBusinessRepository) {
				business := &entities.Business{
					ID:    businessID,
					Name:  "Smith Dental",
					Type:  "dentist",
					Phone: "+1234567890",
				}
				repo.businesses[businessID] = business
			},
			expectedError: false,
			validateResp: func(t *testing.T, resp *dto.BusinessResponse) {
				if resp.Name != "Smith Dental" {
					t.Errorf("expected name unchanged 'Smith Dental', got %s", resp.Name)
				}
				if resp.Phone != "+9876543210" {
					t.Errorf("expected phone '+9876543210', got %s", resp.Phone)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockBusinessRepository{
				businesses: make(map[string]*entities.Business),
			}
			tt.setupMocks(repo)

			service := NewBusinessService(repo, log)

			response, err := service.UpdateBusiness(context.Background(), tt.businessID, tt.request)

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
