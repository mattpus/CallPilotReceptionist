package database

import (
	"context"
	"time"

	"github.com/CallPilotReceptionist/internal/domain/entities"
)

// BusinessRepository defines the interface for business data operations
type BusinessRepository interface {
	Create(ctx context.Context, business *entities.Business) error
	GetByID(ctx context.Context, id string) (*entities.Business, error)
	GetByPhone(ctx context.Context, phone string) (*entities.Business, error)
	Update(ctx context.Context, business *entities.Business) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*entities.Business, error)
}

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id string) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	GetByBusinessID(ctx context.Context, businessID string) ([]*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id string) error
}

// CallRepository defines the interface for call data operations
type CallRepository interface {
	Create(ctx context.Context, call *entities.Call) error
	GetByID(ctx context.Context, id string) (*entities.Call, error)
	GetByProviderCallID(ctx context.Context, providerCallID string) (*entities.Call, error)
	GetByBusinessID(ctx context.Context, businessID string, limit, offset int) ([]*entities.Call, error)
	Update(ctx context.Context, call *entities.Call) error
	Delete(ctx context.Context, id string) error
	GetByDateRange(ctx context.Context, businessID string, startDate, endDate time.Time) ([]*entities.Call, error)
	GetStats(ctx context.Context, businessID string, startDate, endDate time.Time) (*CallStats, error)
}

// InteractionRepository defines the interface for interaction data operations
type InteractionRepository interface {
	Create(ctx context.Context, interaction *entities.Interaction) error
	GetByID(ctx context.Context, id string) (*entities.Interaction, error)
	GetByCallID(ctx context.Context, callID string) ([]*entities.Interaction, error)
	List(ctx context.Context, businessID string, limit, offset int) ([]*entities.Interaction, error)
	Delete(ctx context.Context, id string) error
}

// TranscriptRepository defines the interface for transcript data operations
type TranscriptRepository interface {
	Create(ctx context.Context, transcript *entities.Transcript) error
	CreateBatch(ctx context.Context, transcripts []*entities.Transcript) error
	GetByCallID(ctx context.Context, callID string) ([]*entities.Transcript, error)
	Delete(ctx context.Context, id string) error
	DeleteByCallID(ctx context.Context, callID string) error
}

// AppointmentRepository defines the interface for appointment data operations
type AppointmentRepository interface {
	Create(ctx context.Context, appointment *entities.AppointmentRequest) error
	GetByID(ctx context.Context, id string) (*entities.AppointmentRequest, error)
	GetByCallID(ctx context.Context, callID string) (*entities.AppointmentRequest, error)
	GetByBusinessID(ctx context.Context, businessID string, limit, offset int) ([]*entities.AppointmentRequest, error)
	GetPendingAppointments(ctx context.Context, businessID string) ([]*entities.AppointmentRequest, error)
	Update(ctx context.Context, appointment *entities.AppointmentRequest) error
	Delete(ctx context.Context, id string) error
}

// CallStats represents aggregated call statistics
type CallStats struct {
	TotalCalls       int     `json:"total_calls"`
	CompletedCalls   int     `json:"completed_calls"`
	FailedCalls      int     `json:"failed_calls"`
	TotalDuration    int     `json:"total_duration"` // seconds
	AverageDuration  float64 `json:"average_duration"` // seconds
	TotalCost        float64 `json:"total_cost"`
}
