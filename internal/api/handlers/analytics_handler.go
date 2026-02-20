package handlers

import (
	"net/http"
	"strconv"

	"github.com/CallPilotReceptionist/internal/application/services"
	"github.com/CallPilotReceptionist/internal/infrastructure/http/middleware"
	"github.com/CallPilotReceptionist/pkg/logger"
)

type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
	logger           *logger.Logger
}

func NewAnalyticsHandler(analyticsService *services.AnalyticsService, log *logger.Logger) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
		logger:           log,
	}
}

// GetOverview handles GET /api/v1/analytics/overview
func (h *AnalyticsHandler) GetOverview(w http.ResponseWriter, r *http.Request) {
	businessID := middleware.GetBusinessID(r.Context())

	daysStr := r.URL.Query().Get("days")
	days := 30
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil {
			days = d
		}
	}

	response, err := h.analyticsService.GetOverview(r.Context(), businessID, days)
	if err != nil {
		middleware.RespondError(w, err, h.logger)
		return
	}

	middleware.RespondJSON(w, http.StatusOK, response)
}

// GetCallVolume handles GET /api/v1/analytics/calls
func (h *AnalyticsHandler) GetCallVolume(w http.ResponseWriter, r *http.Request) {
	businessID := middleware.GetBusinessID(r.Context())

	daysStr := r.URL.Query().Get("days")
	days := 30
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil {
			days = d
		}
	}

	response, err := h.analyticsService.GetCallVolume(r.Context(), businessID, days)
	if err != nil {
		middleware.RespondError(w, err, h.logger)
		return
	}

	middleware.RespondJSON(w, http.StatusOK, response)
}
