package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/CallPilotReceptionist/internal/application/dto"
	"github.com/CallPilotReceptionist/internal/domain/errors"
	"github.com/CallPilotReceptionist/pkg/logger"
)

type ErrorMiddleware struct {
	logger *logger.Logger
}

func NewErrorMiddleware(log *logger.Logger) *ErrorMiddleware {
	return &ErrorMiddleware{logger: log}
}

func (m *ErrorMiddleware) Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				m.logger.Error("Panic recovered", nil, map[string]interface{}{
					"error": err,
					"path":  r.URL.Path,
				})

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(dto.ErrorResponse{
					Code:    "INTERNAL_ERROR",
					Message: "Internal server error",
				})
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// RespondError sends a JSON error response based on domain error type
func RespondError(w http.ResponseWriter, err error, log *logger.Logger) {
	w.Header().Set("Content-Type", "application/json")

	var statusCode int
	var response dto.ErrorResponse

	if domainErr, ok := err.(*errors.DomainError); ok {
		response = dto.ErrorResponse{
			Code:    domainErr.Code,
			Message: domainErr.Message,
		}

		switch domainErr.Code {
		case errors.ErrCodeNotFound:
			statusCode = http.StatusNotFound
		case errors.ErrCodeAlreadyExists:
			statusCode = http.StatusConflict
		case errors.ErrCodeInvalidInput, errors.ErrCodeValidationError:
			statusCode = http.StatusBadRequest
		case errors.ErrCodeUnauthorized:
			statusCode = http.StatusUnauthorized
		case errors.ErrCodeForbidden:
			statusCode = http.StatusForbidden
		case errors.ErrCodeProviderError:
			statusCode = http.StatusBadGateway
		default:
			statusCode = http.StatusInternalServerError
		}

		if domainErr.Err != nil {
			response.Details = domainErr.Err.Error()
		}
	} else {
		statusCode = http.StatusInternalServerError
		response = dto.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Internal server error",
		}
		log.Error("Unhandled error", err, nil)
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// RespondJSON sends a JSON success response
func RespondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
