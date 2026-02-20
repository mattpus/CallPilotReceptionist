package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/CallPilotReceptionist/internal/application/dto"
	"github.com/CallPilotReceptionist/internal/application/services"
	"github.com/CallPilotReceptionist/pkg/logger"
)

type contextKey string

const (
	UserIDKey     contextKey = "user_id"
	BusinessIDKey contextKey = "business_id"
	EmailKey      contextKey = "email"
	RoleKey       contextKey = "role"
)

type AuthMiddleware struct {
	authService *services.AuthService
	logger      *logger.Logger
}

func NewAuthMiddleware(authService *services.AuthService, log *logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		logger:      log,
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.respondError(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		// Extract Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			m.respondError(w, http.StatusUnauthorized, "invalid authorization header format")
			return
		}

		token := parts[1]

		// Validate token
		claims, err := m.authService.ValidateToken(token)
		if err != nil {
			m.logger.Warn("Invalid token", map[string]interface{}{
				"error": err.Error(),
			})
			m.respondError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		// Add claims to context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, BusinessIDKey, claims.BusinessID)
		ctx = context.WithValue(ctx, EmailKey, claims.Email)
		ctx = context.WithValue(ctx, RoleKey, claims.Role)

		// Call next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := dto.ErrorResponse{
		Code:    "UNAUTHORIZED",
		Message: message,
	}
	json.NewEncoder(w).Encode(response)
}

// Helper functions to extract values from context
func GetUserID(ctx context.Context) string {
	if val := ctx.Value(UserIDKey); val != nil {
		return val.(string)
	}
	return ""
}

func GetBusinessID(ctx context.Context) string {
	if val := ctx.Value(BusinessIDKey); val != nil {
		return val.(string)
	}
	return ""
}

func GetEmail(ctx context.Context) string {
	if val := ctx.Value(EmailKey); val != nil {
		return val.(string)
	}
	return ""
}

func GetRole(ctx context.Context) string {
	if val := ctx.Value(RoleKey); val != nil {
		return val.(string)
	}
	return ""
}
