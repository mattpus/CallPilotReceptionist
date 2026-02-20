package providers

import (
	"fmt"

	"github.com/CallPilotReceptionist/internal/domain/providers"
	"github.com/CallPilotReceptionist/internal/infrastructure/providers/vapi"
	"github.com/CallPilotReceptionist/pkg/config"
)

type ProviderType string

const (
	ProviderTypeVapi ProviderType = "vapi"
	// Future providers can be added here
	// ProviderTypeTwilio ProviderType = "twilio"
	// ProviderTypeCustom ProviderType = "custom"
)

// ProviderFactory creates voice provider instances based on configuration
type ProviderFactory struct {
	config *config.Config
}

func NewProviderFactory(cfg *config.Config) *ProviderFactory {
	return &ProviderFactory{
		config: cfg,
	}
}

// CreateProvider creates a voice provider instance
// This is the central point for switching between providers
func (f *ProviderFactory) CreateProvider(providerType ProviderType) (providers.VoiceProvider, error) {
	switch providerType {
	case ProviderTypeVapi:
		return vapi.NewVapiProvider(
			f.config.Vapi.APIKey,
			f.config.Vapi.APIBaseURL,
			f.config.Vapi.WebhookURL,
		), nil
	// Future provider implementations:
	// case ProviderTypeTwilio:
	//     return twilio.NewTwilioProvider(...), nil
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", providerType)
	}
}

// GetDefaultProvider returns the default provider configured in the system
func (f *ProviderFactory) GetDefaultProvider() (providers.VoiceProvider, error) {
	// For now, Vapi is the default. This can be made configurable
	return f.CreateProvider(ProviderTypeVapi)
}
