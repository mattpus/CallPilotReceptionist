package services

import (
	"context"
	"testing"
	"time"

	"github.com/CallPilotReceptionist/internal/application/dto"
	"github.com/CallPilotReceptionist/internal/domain/entities"
	"github.com/CallPilotReceptionist/pkg/config"
	"github.com/CallPilotReceptionist/pkg/logger"
)

// Mock repositories
type mockUserRepository struct {
	users map[string]*entities.User
}

func (m *mockUserRepository) Create(ctx context.Context, user *entities.User) error {
	// Generate ID if not set (simulating database behavior)
	if user.ID == "" {
		user.ID = "user-" + time.Now().Format("20060102150405")
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	if user, ok := m.users[id]; ok {
		return user, nil
	}
	return nil, nil
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepository) GetByBusinessID(ctx context.Context, businessID string) ([]*entities.User, error) {
	return nil, nil
}

func (m *mockUserRepository) Update(ctx context.Context, user *entities.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepository) Delete(ctx context.Context, id string) error {
	delete(m.users, id)
	return nil
}

type mockBusinessRepository struct {
	businesses map[string]*entities.Business
}

func (m *mockBusinessRepository) Create(ctx context.Context, business *entities.Business) error {
	// Generate ID if not set (simulating database behavior)
	if business.ID == "" {
		business.ID = "business-" + time.Now().Format("20060102150405")
	}
	m.businesses[business.ID] = business
	return nil
}

func (m *mockBusinessRepository) GetByID(ctx context.Context, id string) (*entities.Business, error) {
	if business, ok := m.businesses[id]; ok {
		return business, nil
	}
	return nil, nil
}

func (m *mockBusinessRepository) GetByPhone(ctx context.Context, phone string) (*entities.Business, error) {
	return nil, nil
}

func (m *mockBusinessRepository) Update(ctx context.Context, business *entities.Business) error {
	m.businesses[business.ID] = business
	return nil
}

func (m *mockBusinessRepository) Delete(ctx context.Context, id string) error {
	delete(m.businesses, id)
	return nil
}

func (m *mockBusinessRepository) List(ctx context.Context, limit, offset int) ([]*entities.Business, error) {
	return nil, nil
}

func newMockAuthService() *AuthService {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			SecretKey:            "test-secret-key",
			AccessTokenDuration:  15 * time.Minute,
			RefreshTokenDuration: 7 * 24 * time.Hour,
		},
	}

	log := logger.New("info", "console")

	return NewAuthService(
		&mockUserRepository{users: make(map[string]*entities.User)},
		&mockBusinessRepository{businesses: make(map[string]*entities.Business)},
		cfg,
		log,
	)
}

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name    string
		req     dto.RegisterRequest
		wantErr bool
	}{
		{
			name: "successful registration",
			req: dto.RegisterRequest{
				BusinessName: "Test Business",
				BusinessType: "dentist",
				Phone:        "+1234567890",
				Email:        "test@example.com",
				Password:     "password123",
			},
			wantErr: false,
		},
		{
			name: "missing email",
			req: dto.RegisterRequest{
				BusinessName: "Test Business",
				BusinessType: "dentist",
				Phone:        "+1234567890",
				Email:        "",
				Password:     "password123",
			},
			wantErr: true,
		},
		{
			name: "missing password",
			req: dto.RegisterRequest{
				BusinessName: "Test Business",
				BusinessType: "dentist",
				Phone:        "+1234567890",
				Email:        "test@example.com",
				Password:     "",
			},
			wantErr: true,
		},
		{
			name: "missing business name",
			req: dto.RegisterRequest{
				BusinessName: "",
				BusinessType: "dentist",
				Phone:        "+1234567890",
				Email:        "test@example.com",
				Password:     "password123",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newMockAuthService()
			ctx := context.Background()

			response, err := service.Register(ctx, tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if response == nil {
					t.Error("Register() expected response, got nil")
					return
				}
				if response.AccessToken == "" {
					t.Error("Register() expected access token, got empty")
				}
				if response.RefreshToken == "" {
					t.Error("Register() expected refresh token, got empty")
				}
				if response.User.Email != tt.req.Email {
					t.Errorf("Register() user email = %v, want %v", response.User.Email, tt.req.Email)
				}
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	service := newMockAuthService()
	ctx := context.Background()

	// First register a user
	registerReq := dto.RegisterRequest{
		BusinessName: "Test Business",
		BusinessType: "dentist",
		Phone:        "+1234567890",
		Email:        "test@example.com",
		Password:     "password123",
	}
	_, err := service.Register(ctx, registerReq)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	tests := []struct {
		name    string
		req     dto.LoginRequest
		wantErr bool
	}{
		{
			name: "successful login",
			req: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "wrong password",
			req: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			wantErr: true,
		},
		{
			name: "non-existent user",
			req: dto.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "missing email",
			req: dto.LoginRequest{
				Email:    "",
				Password: "password123",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := service.Login(ctx, tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if response == nil {
					t.Error("Login() expected response, got nil")
					return
				}
				if response.AccessToken == "" {
					t.Error("Login() expected access token, got empty")
				}
				if response.RefreshToken == "" {
					t.Error("Login() expected refresh token, got empty")
				}
			}
		})
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	service := newMockAuthService()
	ctx := context.Background()

	// Register and get token
	registerReq := dto.RegisterRequest{
		BusinessName: "Test Business",
		BusinessType: "dentist",
		Phone:        "+1234567890",
		Email:        "test@example.com",
		Password:     "password123",
	}
	response, err := service.Register(ctx, registerReq)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   response.AccessToken,
			wantErr: false,
		},
		{
			name:    "invalid token",
			token:   "invalid.token.here",
			wantErr: true,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := service.ValidateToken(tt.token)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if claims == nil {
					t.Error("ValidateToken() expected claims, got nil")
					return
				}
				if claims.Email != "test@example.com" {
					t.Errorf("ValidateToken() email = %v, want test@example.com", claims.Email)
				}
			}
		})
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	service := newMockAuthService()
	ctx := context.Background()

	// Register and get tokens
	registerReq := dto.RegisterRequest{
		BusinessName: "Test Business",
		BusinessType: "dentist",
		Phone:        "+1234567890",
		Email:        "test@example.com",
		Password:     "password123",
	}
	response, err := service.Register(ctx, registerReq)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid refresh token",
			token:   response.RefreshToken,
			wantErr: false,
		},
		{
			name:    "invalid refresh token",
			token:   "invalid.token.here",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			refreshResponse, err := service.RefreshToken(ctx, tt.token)

			if (err != nil) != tt.wantErr {
				t.Errorf("RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if refreshResponse == nil {
					t.Error("RefreshToken() expected response, got nil")
					return
				}
				if refreshResponse.AccessToken == "" {
					t.Error("RefreshToken() expected access token, got empty")
				}
			}
		})
	}
}
