package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/CallPilotReceptionist/internal/application/dto"
	"github.com/CallPilotReceptionist/internal/application/services"
	"github.com/CallPilotReceptionist/internal/infrastructure/http/middleware"
	"github.com/CallPilotReceptionist/pkg/logger"
)

type AuthHandler struct {
	authService *services.AuthService
	logger      *logger.Logger
}

func NewAuthHandler(authService *services.AuthService, log *logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      log,
	}
}

// Register handles POST /api/v1/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, err, h.logger)
		return
	}

	response, err := h.authService.Register(r.Context(), req)
	if err != nil {
		middleware.RespondError(w, err, h.logger)
		return
	}

	middleware.RespondJSON(w, http.StatusCreated, response)
}

// Login handles POST /api/v1/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, err, h.logger)
		return
	}

	response, err := h.authService.Login(r.Context(), req)
	if err != nil {
		middleware.RespondError(w, err, h.logger)
		return
	}

	middleware.RespondJSON(w, http.StatusOK, response)
}

// RefreshToken handles POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, err, h.logger)
		return
	}

	response, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		middleware.RespondError(w, err, h.logger)
		return
	}

	middleware.RespondJSON(w, http.StatusOK, response)
}

// Logout handles POST /api/v1/auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// In a stateless JWT system, logout is typically handled client-side
	// by deleting the token. We can add token blacklisting in the future.
	middleware.RespondJSON(w, http.StatusOK, dto.SuccessResponse{
		Message: "Logged out successfully",
	})
}
