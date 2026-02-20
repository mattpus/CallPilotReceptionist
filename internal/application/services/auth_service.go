package services

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/CallPilotReceptionist/internal/application/dto"
	"github.com/CallPilotReceptionist/internal/domain/entities"
	"github.com/CallPilotReceptionist/internal/domain/errors"
	"github.com/CallPilotReceptionist/internal/infrastructure/database"
	"github.com/CallPilotReceptionist/pkg/config"
	"github.com/CallPilotReceptionist/pkg/logger"
)

type AuthService struct {
	userRepo     database.UserRepository
	businessRepo database.BusinessRepository
	config       *config.Config
	logger       *logger.Logger
}

func NewAuthService(
	userRepo database.UserRepository,
	businessRepo database.BusinessRepository,
	cfg *config.Config,
	log *logger.Logger,
) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		businessRepo: businessRepo,
		config:       cfg,
		logger:       log,
	}
}

type Claims struct {
	UserID     string `json:"user_id"`
	BusinessID string `json:"business_id"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	jwt.RegisteredClaims
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.LoginResponse, error) {
	// Validate input
	if req.Email == "" || req.Password == "" || req.BusinessName == "" || req.Phone == "" {
		return nil, errors.NewValidationError("email, password, business name, and phone are required")
	}

	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.NewAlreadyExistsError("user", "email", req.Email)
	}

	// Create business
	business, err := entities.NewBusiness(req.BusinessName, req.BusinessType, req.Phone, nil)
	if err != nil {
		return nil, err
	}

	if err := s.businessRepo.Create(ctx, business); err != nil {
		s.logger.Error("Failed to create business", err, map[string]interface{}{
			"business_name": req.BusinessName,
		})
		return nil, err
	}

	// Hash password
	hashedPassword, err := database.HashPassword(req.Password)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	// Create user
	user, err := entities.NewUser(business.ID, req.Email, hashedPassword, entities.UserRoleOwner)
	if err != nil {
		return nil, err
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error("Failed to create user", err, map[string]interface{}{
			"email": req.Email,
		})
		return nil, err
	}

	s.logger.Info("User registered successfully", map[string]interface{}{
		"user_id":     user.ID,
		"business_id": business.ID,
		"email":       user.Email,
	})

	// Generate tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:         user.ID,
			BusinessID: user.BusinessID,
			Email:      user.Email,
			Role:       string(user.Role),
			CreatedAt:  user.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Validate input
	if req.Email == "" || req.Password == "" {
		return nil, errors.NewValidationError("email and password are required")
	}

	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.NewUnauthorizedError("invalid credentials")
	}
	
	// Check if user exists
	if user == nil {
		s.logger.Warn("Login attempt for non-existent user", map[string]interface{}{
			"email": req.Email,
		})
		return nil, errors.NewUnauthorizedError("invalid credentials")
	}

	// Compare password
	if err := database.ComparePassword(user.PasswordHash, req.Password); err != nil {
		s.logger.Warn("Failed login attempt", map[string]interface{}{
			"email": req.Email,
		})
		return nil, errors.NewUnauthorizedError("invalid credentials")
	}

	s.logger.Info("User logged in successfully", map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
	})

	// Generate tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:         user.ID,
			BusinessID: user.BusinessID,
			Email:      user.Email,
			Role:       string(user.Role),
			CreatedAt:  user.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshTokenResponse, error) {
	// Parse and validate refresh token
	claims, err := s.validateToken(refreshToken)
	if err != nil {
		return nil, errors.NewUnauthorizedError("invalid refresh token")
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.NewUnauthorizedError("user not found")
	}

	// Generate new access token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	return &dto.RefreshTokenResponse{
		AccessToken: accessToken,
	}, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	return s.validateToken(tokenString)
}

func (s *AuthService) generateAccessToken(user *entities.User) (string, error) {
	claims := &Claims{
		UserID:     user.ID,
		BusinessID: user.BusinessID,
		Email:      user.Email,
		Role:       string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.JWT.AccessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "vapi-integration",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.SecretKey))
}

func (s *AuthService) generateRefreshToken(user *entities.User) (string, error) {
	claims := &Claims{
		UserID:     user.ID,
		BusinessID: user.BusinessID,
		Email:      user.Email,
		Role:       string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.JWT.RefreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "vapi-integration",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.SecretKey))
}

func (s *AuthService) validateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWT.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.NewUnauthorizedError("invalid token")
}
