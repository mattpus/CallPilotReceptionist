package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/CallPilotReceptionist/internal/application/dto"
	"github.com/CallPilotReceptionist/internal/application/services"
	"github.com/CallPilotReceptionist/internal/infrastructure/http/middleware"
	"github.com/CallPilotReceptionist/pkg/logger"
)

type CallHandler struct {
	callService *services.CallService
	logger      *logger.Logger
}

func NewCallHandler(callService *services.CallService, log *logger.Logger) *CallHandler {
	return &CallHandler{
		callService: callService,
		logger:      log,
	}
}

// InitiateCall handles POST /api/v1/calls
func (h *CallHandler) InitiateCall(w http.ResponseWriter, r *http.Request) {
	businessID := middleware.GetBusinessID(r.Context())

	var req dto.InitiateCallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, err, h.logger)
		return
	}

	response, err := h.callService.InitiateCall(r.Context(), businessID, req)
	if err != nil {
		middleware.RespondError(w, err, h.logger)
		return
	}

	middleware.RespondJSON(w, http.StatusCreated, response)
}

// ListCalls handles GET /api/v1/calls
func (h *CallHandler) ListCalls(w http.ResponseWriter, r *http.Request) {
	businessID := middleware.GetBusinessID(r.Context())

	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	status := r.URL.Query().Get("status")

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

	req := dto.ListCallsRequest{
		Limit:  limit,
		Offset: offset,
		Status: status,
	}

	response, err := h.callService.ListCalls(r.Context(), businessID, req)
	if err != nil {
		middleware.RespondError(w, err, h.logger)
		return
	}

	middleware.RespondJSON(w, http.StatusOK, response)
}

// GetCall handles GET /api/v1/calls/:id
func (h *CallHandler) GetCall(w http.ResponseWriter, r *http.Request) {
	businessID := middleware.GetBusinessID(r.Context())
	vars := mux.Vars(r)
	callID := vars["id"]

	response, err := h.callService.GetCall(r.Context(), businessID, callID)
	if err != nil {
		middleware.RespondError(w, err, h.logger)
		return
	}

	middleware.RespondJSON(w, http.StatusOK, response)
}

// GetTranscript handles GET /api/v1/calls/:id/transcript
func (h *CallHandler) GetTranscript(w http.ResponseWriter, r *http.Request) {
	businessID := middleware.GetBusinessID(r.Context())
	vars := mux.Vars(r)
	callID := vars["id"]

	response, err := h.callService.GetTranscript(r.Context(), businessID, callID)
	if err != nil {
		middleware.RespondError(w, err, h.logger)
		return
	}

	middleware.RespondJSON(w, http.StatusOK, response)
}

// HandleWebhook handles POST /api/v1/webhooks/vapi
func (h *CallHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	// Read raw body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("Failed to read webhook body", err, nil)
		middleware.RespondError(w, err, h.logger)
		return
	}

	// Get signature from header
	signature := r.Header.Get("X-Vapi-Signature")

	// Process webhook
	if err := h.callService.HandleWebhook(r.Context(), body, signature); err != nil {
		middleware.RespondError(w, err, h.logger)
		return
	}

	middleware.RespondJSON(w, http.StatusOK, dto.SuccessResponse{
		Message: "Webhook processed successfully",
	})
}
