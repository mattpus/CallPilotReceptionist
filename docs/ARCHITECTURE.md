# Vapi AI Integration Backend - Architecture Documentation

## Table of Contents
1. [Overview](#overview)
2. [Architecture Diagrams](#architecture-diagrams)
3. [System Architecture](#system-architecture)
4. [Provider Abstraction](#provider-abstraction)
5. [Authentication Flow](#authentication-flow)
6. [Call Flow](#call-flow)
7. [Database Schema](#database-schema)
8. [Deployment Architecture](#deployment-architecture)
9. [Testing Strategy](#testing-strategy)
10. [Provider Switching Guide](#provider-switching-guide)

---

## Overview

This document provides comprehensive architectural documentation for the Vapi AI Integration Backend. The system is designed to provide AI voice assistants for small businesses (e.g., dentists, salons) with a focus on:

- **Provider Abstraction**: Easy switching between voice AI providers (Vapi AI, Twilio, etc.)
- **Hexagonal Architecture**: Clean separation of concerns with domain, application, infrastructure, and API layers
- **Comprehensive Testing**: Unit tests, integration tests, and mocks for all components
- **Scalability**: Horizontal scaling with multiple application instances
- **Security**: JWT authentication, webhook signature validation, password hashing

---

## Architecture Diagrams

All diagrams are available in PlantUML format in the `docs/diagrams/` directory:

- `architecture.puml` - System architecture overview
- `provider-abstraction.puml` - Provider abstraction pattern
- `auth-flow-sequence.puml` - Authentication and authorization flow
- `call-flow-sequence.puml` - Call initiation and webhook processing
- `database-schema.puml` - Database schema and relationships
- `deployment.puml` - Deployment architecture

### Generating Diagrams

To generate PNG/SVG images from PlantUML files:

```bash
# Install PlantUML (requires Java)
brew install plantuml  # macOS
apt-get install plantuml  # Ubuntu/Debian

# Generate all diagrams
cd docs/diagrams
for file in *.puml; do
    plantuml "$file"
done

# This will create PNG files for each diagram
```

Alternatively, use online PlantUML editors:
- https://www.plantuml.com/plantuml/uml/
- https://plantuml-editor.kkeisuke.com/

---

## System Architecture

### Hexagonal Architecture (Ports & Adapters)

The system follows hexagonal architecture principles with clear separation between:

1. **API Layer** (Adapters)
   - HTTP handlers for REST endpoints
   - Request/response DTOs
   - Middleware (auth, logging, CORS, error handling)

2. **Application Layer** (Use Cases)
   - Business logic services
   - Service orchestration
   - Transaction management

3. **Domain Layer** (Core)
   - Business entities
   - Domain errors
   - Provider interfaces (ports)

4. **Infrastructure Layer** (Adapters)
   - Database repositories
   - Voice provider implementations
   - External API clients

### Key Benefits

- **Testability**: Mock external dependencies easily
- **Maintainability**: Changes isolated to specific layers
- **Flexibility**: Swap implementations without changing business logic
- **Clarity**: Clear dependency direction (inward toward domain)

### Layer Dependencies

```
API Layer → Application Layer → Domain Layer
                ↓                    ↑
         Infrastructure Layer ------┘
```

**Rule**: Dependencies always point inward. Infrastructure depends on domain (interfaces), never the reverse.

---

## Provider Abstraction

### The VoiceProvider Interface

The core abstraction that enables easy provider switching:

```go
type VoiceProvider interface {
    InitiateCall(ctx context.Context, request InitiateCallRequest) (*CallSession, error)
    HandleWebhook(ctx context.Context, payload []byte, signature string) (*CallEvent, error)
    GetCallDetails(ctx context.Context, callID string) (*CallDetails, error)
    GetTranscript(ctx context.Context, callID string) ([]Message, error)
    UpdateCall(ctx context.Context, callID string, updates map[string]interface{}) error
    GetCallRecording(ctx context.Context, callID string) (*Recording, error)
    ListCalls(ctx context.Context, businessID string, filters ListCallsFilters) ([]CallSummary, error)
    CancelCall(ctx context.Context, callID string) error
    GetProviderName() string
}
```

### Implementation Strategy

1. **Vapi AI Implementation** (`internal/infrastructure/providers/vapi/vapi_provider.go`)
   - HTTP API client for Vapi AI
   - Webhook signature validation (HMAC SHA256)
   - Status code mapping
   - Error handling and retries

2. **Factory Pattern** (`internal/infrastructure/providers/factory.go`)
   - Creates provider based on configuration
   - Centralizes provider instantiation
   - Easy to add new providers

### Switching Providers

To switch from Vapi AI to another provider (e.g., Twilio):

**Step 1**: Implement the `VoiceProvider` interface

```go
// internal/infrastructure/providers/twilio/twilio_provider.go
type TwilioProvider struct {
    accountSID string
    authToken  string
    httpClient *http.Client
}

func (tp *TwilioProvider) InitiateCall(ctx context.Context, req InitiateCallRequest) (*CallSession, error) {
    // Twilio-specific implementation
}

// ... implement all other methods
```

**Step 2**: Add to factory

```go
// internal/infrastructure/providers/factory.go
func (f *ProviderFactory) NewVoiceProvider(providerType string) (providers.VoiceProvider, error) {
    switch providerType {
    case "vapi":
        return f.createVapiProvider()
    case "twilio":
        return f.createTwilioProvider()  // Add this
    default:
        return nil, errors.New("unknown provider")
    }
}
```

**Step 3**: Update configuration

```bash
# .env
VOICE_PROVIDER=twilio
TWILIO_ACCOUNT_SID=your_account_sid
TWILIO_AUTH_TOKEN=your_auth_token
```

**That's it!** No changes needed to:
- CallService
- HTTP handlers
- Domain entities
- Database schema
- API endpoints

---

## Authentication Flow

### JWT Token Strategy

- **Access Tokens**: 15-minute expiry, used for API authentication
- **Refresh Tokens**: 7-day expiry, used to obtain new access tokens
- **Claims**: userID, businessID, email, role

### Registration Flow

1. User submits registration form (business name, email, password)
2. AuthService validates input (email format, password strength)
3. Check if email already exists
4. Hash password with bcrypt (cost=10)
5. Create business record in database
6. Create user record linked to business
7. Generate JWT tokens (access + refresh)
8. Return tokens to client

### Login Flow

1. User submits credentials (email, password)
2. AuthService retrieves user by email
3. Compare password with bcrypt.CompareHashAndPassword
4. If match: generate JWT tokens
5. If mismatch: return 401 Unauthorized
6. Return tokens to client

### Protected Endpoint Access

1. Client sends request with `Authorization: Bearer <token>`
2. Auth middleware extracts token from header
3. Validate token (signature, expiry)
4. If valid: extract claims and store in context
5. If invalid/expired: return 401 Unauthorized
6. Handler accesses claims via context helpers (GetUserID, GetBusinessID, etc.)

### Token Refresh

1. Client sends refresh token to `/api/v1/auth/refresh`
2. Validate refresh token (signature, expiry)
3. Retrieve user from database
4. Generate new access token
5. Return new access token to client

---

## Call Flow

### Call Initiation

1. **User Request**: Business user initiates call via dashboard
   - POST `/api/v1/calls` with phone number and config

2. **Service Layer**:
   - Create call record in database (status: "initiated")
   - Call VoiceProvider.InitiateCall()
   - Provider makes HTTP request to Vapi AI API
   - Update call with provider_call_id

3. **Response**:
   - Return call details to user
   - Call status: "initiated"

### Webhook Processing

1. **Vapi AI Event**: Vapi sends webhook to `/api/v1/webhooks/vapi`
   - Event types: call.started, call.ended, call.failed
   - Includes X-Vapi-Signature header

2. **Signature Validation**:
   - CallService delegates to VoiceProvider.HandleWebhook()
   - Provider validates HMAC SHA256 signature
   - If invalid: return 401 Unauthorized

3. **Database Update**:
   - Retrieve call by provider_call_id
   - Update call status based on event
   - Update timestamps (started_at, ended_at)

4. **Async Transcript Fetch** (on call.ended):
   - Start goroutine to fetch transcript
   - Call VoiceProvider.GetTranscript()
   - Store transcript messages in database
   - Parse interactions (appointments, questions, etc.)

### Call Transcript Retrieval

1. **User Request**: GET `/api/v1/calls/{id}/transcript`
2. **Service Layer**:
   - Verify user owns call (businessID match)
   - Retrieve transcript from database
   - Order by timestamp
3. **Response**: Array of transcript messages with role (assistant/user)

---

## Database Schema

### Tables

1. **businesses**
   - Primary business information
   - JSONB settings for flexible configuration
   - Indexed: id, phone

2. **users**
   - User accounts linked to businesses
   - Password hash (bcrypt)
   - Role-based access control
   - Indexed: id, email, business_id

3. **calls**
   - Call records with provider link
   - Status tracking (initiated, in_progress, completed, failed)
   - Duration and cost tracking
   - Indexed: id, business_id, provider_call_id, status, created_at

4. **interactions**
   - Parsed interactions from calls
   - JSONB content for flexibility
   - Types: appointment_request, question, complaint, etc.
   - Indexed: id, call_id, type

5. **transcripts**
   - Conversation messages from calls
   - Role: assistant or user
   - Timestamp for ordering
   - Indexed: id, call_id, timestamp

6. **appointments**
   - Appointment requests extracted from calls
   - Status workflow: pending → confirmed → completed/cancelled
   - Indexed: id, call_id, business_id, status, requested_date

### Relationships

- businesses 1:N users
- businesses 1:N calls
- businesses 1:N appointments
- calls 1:N interactions
- calls 1:N transcripts
- calls 1:1 appointments (may generate)

### Indexes

All foreign keys are indexed for query performance. Additional indexes on:
- Timestamp columns for sorting
- Status columns for filtering
- Phone numbers for lookups

---

## Deployment Architecture

### Docker Setup

**Multi-Container Architecture**:
- Load Balancer (Nginx/Traefik)
- Go Application (3 replicas for horizontal scaling)
- Supabase (PostgreSQL)
- Redis (optional, for caching)
- Monitoring (Prometheus + Grafana)

### Environment Configuration

```bash
# Database
DATABASE_URL=postgresql://user:pass@supabase:5432/dbname

# Voice Provider
VOICE_PROVIDER=vapi
VAPI_API_KEY=your_api_key
VAPI_WEBHOOK_SECRET=your_webhook_secret

# JWT
JWT_SECRET=your_jwt_secret
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

# Server
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
SERVER_SHUTDOWN_TIMEOUT=30s

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://yourdomain.com
```

### Scaling Strategy

1. **Horizontal Scaling**: Add more application instances
2. **Database Connection Pooling**: Supabase handles connection pooling
3. **Stateless Design**: JWT tokens, no server-side sessions
4. **Async Processing**: Transcript fetching in goroutines
5. **Health Checks**: `/health` endpoint for load balancer

### Monitoring

**Metrics to Track**:
- Active calls count
- API endpoint latency (p50, p95, p99)
- Error rate by endpoint
- Database query performance
- Provider API latency
- Webhook processing time

**Logging**:
- Structured logging with zerolog
- Log levels: debug, info, warn, error
- Contextual logging (request ID, user ID)

---

## Testing Strategy

### Unit Tests

**Coverage Target**: 80%+

**Test Files**:
- `auth_service_test.go` - Authentication logic
- `call_service_test.go` - Call management
- `business_service_test.go` - Business operations
- `entities_test.go` - Domain entity validation

**Approach**:
- Table-driven tests (Go idiom)
- Mock repositories for database isolation
- Mock providers for external service isolation
- Test success and failure paths
- Test edge cases (empty inputs, invalid formats)

### Integration Tests

**Purpose**: Test HTTP handlers with mock services

**Test Files**:
- `auth_handler_test.go`
- `call_handler_test.go`
- `business_handler_test.go`

**Approach**:
- Use httptest.NewRecorder
- Mock services
- Test middleware behavior (auth, logging, error handling)
- Verify HTTP status codes
- Validate response JSON

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test ./internal/application/services/

# Run with verbose output
go test -v ./...

# Run tests in parallel
go test -parallel 4 ./...
```

### Test Organization

```
internal/
├── application/
│   └── services/
│       ├── auth_service.go
│       ├── auth_service_test.go
│       ├── call_service.go
│       └── call_service_test.go
├── domain/
│   └── entities/
│       ├── entities.go
│       └── entities_test.go
└── api/
    └── handlers/
        ├── auth_handler.go
        └── auth_handler_test.go
```

---

## Provider Switching Guide

### Why Provider Abstraction Matters

In the voice AI space, providers may:
- Change pricing models
- Deprecate features
- Experience downtime
- Have regional limitations
- Offer different capabilities

With our abstraction, switching providers requires **zero changes to business logic**.

### Adding a New Provider

**Example**: Adding Twilio as a voice provider

#### Step 1: Create Provider Implementation

```go
// internal/infrastructure/providers/twilio/twilio_provider.go
package twilio

import (
    "context"
    "vapiAIIntegration/internal/domain/providers"
)

type TwilioProvider struct {
    accountSID string
    authToken  string
    httpClient *http.Client
    logger     *zerolog.Logger
}

func NewTwilioProvider(accountSID, authToken string, logger *zerolog.Logger) *TwilioProvider {
    return &TwilioProvider{
        accountSID: accountSID,
        authToken:  authToken,
        httpClient: &http.Client{Timeout: 30 * time.Second},
        logger:     logger,
    }
}

func (tp *TwilioProvider) InitiateCall(ctx context.Context, request providers.InitiateCallRequest) (*providers.CallSession, error) {
    // Implement Twilio-specific call initiation
    // Example: use Twilio's REST API to create a call
    url := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Calls.json", tp.accountSID)
    
    // ... make HTTP request to Twilio
    // ... parse response
    // ... return CallSession
}

// Implement all other interface methods...
```

#### Step 2: Add to Factory

```go
// internal/infrastructure/providers/factory.go

func (f *ProviderFactory) createTwilioProvider() (*twilio.TwilioProvider, error) {
    accountSID := f.config.TwilioAccountSID
    authToken := f.config.TwilioAuthToken
    
    if accountSID == "" || authToken == "" {
        return nil, errors.New("Twilio credentials not configured")
    }
    
    return twilio.NewTwilioProvider(accountSID, authToken, f.logger), nil
}

func (f *ProviderFactory) NewVoiceProvider(providerType string) (providers.VoiceProvider, error) {
    switch providerType {
    case "vapi":
        return f.createVapiProvider()
    case "twilio":
        return f.createTwilioProvider()
    default:
        return nil, fmt.Errorf("unknown provider type: %s", providerType)
    }
}
```

#### Step 3: Update Configuration

```go
// pkg/config/config.go

type Config struct {
    // ... existing fields
    
    // Twilio Configuration
    TwilioAccountSID string `env:"TWILIO_ACCOUNT_SID"`
    TwilioAuthToken  string `env:"TWILIO_AUTH_TOKEN"`
}
```

#### Step 4: Update Environment

```bash
# .env
VOICE_PROVIDER=twilio
TWILIO_ACCOUNT_SID=ACxxxxxxxxxxxxxxxxxxxxx
TWILIO_AUTH_TOKEN=your_auth_token_here
```

#### Step 5: Test

```go
// internal/infrastructure/providers/twilio/twilio_provider_test.go

func TestTwilioProvider_InitiateCall(t *testing.T) {
    provider := NewTwilioProvider("test_sid", "test_token", logger.New("test"))
    
    request := providers.InitiateCallRequest{
        PhoneNumber: "+1234567890",
        BusinessID:  uuid.New(),
    }
    
    session, err := provider.InitiateCall(context.Background(), request)
    
    // ... assertions
}
```

### Provider Comparison Matrix

| Feature | Vapi AI | Twilio | Custom |
|---------|---------|--------|--------|
| Voice Recognition | ✅ | ✅ | ⚠️ |
| Webhook Support | ✅ | ✅ | ✅ |
| Transcription | ✅ | ✅ | ⚠️ |
| Call Recording | ✅ | ✅ | ✅ |
| Pricing Model | Per minute | Per minute | Custom |
| Setup Complexity | Low | Medium | High |

### Migration Checklist

When switching providers:

- [ ] Implement VoiceProvider interface
- [ ] Add to factory
- [ ] Update configuration
- [ ] Test all interface methods
- [ ] Test webhook handling
- [ ] Update environment variables
- [ ] Test call initiation
- [ ] Test transcript retrieval
- [ ] Update documentation
- [ ] Monitor error rates post-migration

---

## Best Practices

### Code Organization

1. **Domain-Driven Design**: Keep business logic in domain layer
2. **Interface Segregation**: Small, focused interfaces
3. **Dependency Injection**: Pass dependencies explicitly
4. **Error Handling**: Use domain errors with codes
5. **Logging**: Structured logging with context

### Security

1. **JWT Tokens**: Short-lived access tokens
2. **Password Hashing**: bcrypt with appropriate cost
3. **Webhook Validation**: Verify signatures
4. **Input Validation**: Validate all user input
5. **CORS**: Configure allowed origins

### Performance

1. **Database Indexes**: Index foreign keys and frequently queried columns
2. **Connection Pooling**: Use Supabase connection pooler
3. **Async Processing**: Use goroutines for long-running tasks
4. **Caching**: Optional Redis for frequently accessed data
5. **Query Optimization**: Use EXPLAIN ANALYZE

### Maintainability

1. **Documentation**: Keep docs up-to-date
2. **Tests**: Maintain high test coverage
3. **Linting**: Use golangci-lint
4. **Code Reviews**: Review all changes
5. **Versioning**: Use semantic versioning

---

## Troubleshooting

### Common Issues

**Problem**: Webhook signature validation fails
- **Solution**: Check VAPI_WEBHOOK_SECRET matches Vapi dashboard
- **Debug**: Log received signature and computed signature

**Problem**: Database connection timeout
- **Solution**: Check DATABASE_URL, verify Supabase is running
- **Debug**: Test connection with psql

**Problem**: JWT token invalid
- **Solution**: Check JWT_SECRET, verify token expiry
- **Debug**: Decode token at jwt.io

**Problem**: Provider API rate limit
- **Solution**: Implement retry with exponential backoff
- **Debug**: Check provider dashboard for rate limits

---

## Additional Resources

- [Go Best Practices](https://golang.org/doc/effective_go)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Vapi AI Documentation](https://docs.vapi.ai/)
- [Supabase Documentation](https://supabase.com/docs)
- [JWT Best Practices](https://tools.ietf.org/html/rfc8725)

---

## Contributing

When contributing to this project:

1. Follow existing code structure and patterns
2. Write tests for all new functionality
3. Update documentation for changes
4. Use conventional commit messages
5. Keep provider abstraction intact

---

## License

[Your License Here]

---

## Support

For questions or issues:
- GitHub Issues: [repository URL]
- Email: [support email]
- Documentation: This file and `/docs` directory
