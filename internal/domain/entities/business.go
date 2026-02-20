package entities

import (
	"time"

	"github.com/CallPilotReceptionist/internal/domain/errors"
)

type Business struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Phone     string                 `json:"phone"`
	Settings  map[string]interface{} `json:"settings"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

func NewBusiness(name, businessType, phone string, settings map[string]interface{}) (*Business, error) {
	if name == "" {
		return nil, errors.NewValidationError("business name is required")
	}
	if phone == "" {
		return nil, errors.NewValidationError("business phone is required")
	}

	now := time.Now()
	return &Business{
		Name:      name,
		Type:      businessType,
		Phone:     phone,
		Settings:  settings,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (b *Business) Update(name, businessType, phone string, settings map[string]interface{}) error {
	if name != "" {
		b.Name = name
	}
	if businessType != "" {
		b.Type = businessType
	}
	if phone != "" {
		b.Phone = phone
	}
	if settings != nil {
		b.Settings = settings
	}
	b.UpdatedAt = time.Now()
	return nil
}

func (b *Business) Validate() error {
	if b.Name == "" {
		return errors.NewValidationError("business name is required")
	}
	if b.Phone == "" {
		return errors.NewValidationError("business phone is required")
	}
	return nil
}
