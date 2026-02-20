package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/CallPilotReceptionist/internal/application/dto"
	"github.com/CallPilotReceptionist/internal/application/services"
	"github.com/CallPilotReceptionist/internal/infrastructure/http/middleware"
	"github.com/CallPilotReceptionist/pkg/logger"
)

type InteractionHandler struct {
	interactionService *services.InteractionService
	logger             *logger.Logger
}

func NewInteractionHandler(interactionService *services.InteractionService, log *logger.Logger) *InteractionHandler {
	return &InteractionHandler{
		interactionService: interactionService,
		logger:             log,
	}
}

// GetCallInteractions handles GET /api/v1/calls/:id/interactions
func (h *InteractionHandler) GetCallInteractions(w http.ResponseWriter, r *http.Request) {
	businessID := middleware.GetBusinessID(r.Context())
	vars := mux.Vars(r)
	callID := vars["id"]

	response, err := h.interactionService.GetCallInteractions(r.Context(), businessID, callID)
	if err != nil {
		middleware.RespondError(w, err, h.logger)
		return
	}

	middleware.RespondJSON(w, http.StatusOK, response)
}

// ListInteractions handles GET /api/v1/interactions
func (h *InteractionHandler) ListInteractions(w http.ResponseWriter, r *http.Request) {
	businessID := middleware.GetBusinessID(r.Context())

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	response, err := h.interactionService.ListInteractions(r.Context(), businessID, limit, offset)
	if err != nil {
		middleware.RespondError(w, err, h.logger)
		return
	}

	middleware.RespondJSON(w, http.StatusOK, response)
}

// ListAppointments handles GET /api/v1/appointments
func (h *InteractionHandler) ListAppointments(w http.ResponseWriter, r *http.Request) {
	businessID := middleware.GetBusinessID(r.Context())

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	response, err := h.interactionService.GetAppointments(r.Context(), businessID, limit, offset)
	if err != nil {
		middleware.RespondError(w, err, h.logger)
		return
	}

	middleware.RespondJSON(w, http.StatusOK, response)
}

// UpdateAppointmentStatus handles PATCH /api/v1/appointments/:id
func (h *InteractionHandler) UpdateAppointmentStatus(w http.ResponseWriter, r *http.Request) {
	businessID := middleware.GetBusinessID(r.Context())
	vars := mux.Vars(r)
	appointmentID := vars["id"]

	var req dto.UpdateAppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, err, h.logger)
		return
	}

	response, err := h.interactionService.UpdateAppointmentStatus(r.Context(), businessID, appointmentID, req)
	if err != nil {
		middleware.RespondError(w, err, h.logger)
		return
	}

	middleware.RespondJSON(w, http.StatusOK, response)
}
