package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/CallPilotReceptionist/internal/application/dto"
	"github.com/CallPilotReceptionist/internal/application/services"
	"github.com/CallPilotReceptionist/internal/infrastructure/http/middleware"
	"github.com/CallPilotReceptionist/pkg/logger"
)

type BusinessHandler struct {
	businessService *services.BusinessService
	logger          *logger.Logger
}

func NewBusinessHandler(businessService *services.BusinessService, log *logger.Logger) *BusinessHandler {
	return &BusinessHandler{
		businessService: businessService,
		logger:          log,
	}
}

// GetBusiness handles GET /api/v1/businesses/me
func (h *BusinessHandler) GetBusiness(w http.ResponseWriter, r *http.Request) {
	businessID := middleware.GetBusinessID(r.Context())

	response, err := h.businessService.GetBusiness(r.Context(), businessID)
	if err != nil {
		middleware.RespondError(w, err, h.logger)
		return
	}

	middleware.RespondJSON(w, http.StatusOK, response)
}

// UpdateBusiness handles PUT /api/v1/businesses/me
func (h *BusinessHandler) UpdateBusiness(w http.ResponseWriter, r *http.Request) {
	businessID := middleware.GetBusinessID(r.Context())

	var req dto.UpdateBusinessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, err, h.logger)
		return
	}

	response, err := h.businessService.UpdateBusiness(r.Context(), businessID, req)
	if err != nil {
		middleware.RespondError(w, err, h.logger)
		return
	}

	middleware.RespondJSON(w, http.StatusOK, response)
}
