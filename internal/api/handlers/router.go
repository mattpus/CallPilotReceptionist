package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/CallPilotReceptionist/internal/application/services"
	"github.com/CallPilotReceptionist/internal/infrastructure/http/middleware"
	"github.com/CallPilotReceptionist/pkg/logger"
)

type Router struct {
	router              *mux.Router
	authMiddleware      *middleware.AuthMiddleware
	loggingMiddleware   *middleware.LoggingMiddleware
	corsMiddleware      *middleware.CORSMiddleware
	errorMiddleware     *middleware.ErrorMiddleware
	authHandler         *AuthHandler
	businessHandler     *BusinessHandler
	callHandler         *CallHandler
	analyticsHandler    *AnalyticsHandler
	interactionHandler  *InteractionHandler
}

func NewRouter(
	authService *services.AuthService,
	businessService *services.BusinessService,
	callService *services.CallService,
	analyticsService *services.AnalyticsService,
	interactionService *services.InteractionService,
	log *logger.Logger,
) *Router {
	r := &Router{
		router:              mux.NewRouter(),
		authMiddleware:      middleware.NewAuthMiddleware(authService, log),
		loggingMiddleware:   middleware.NewLoggingMiddleware(log),
		corsMiddleware:      middleware.NewCORSMiddleware(nil, nil, nil),
		errorMiddleware:     middleware.NewErrorMiddleware(log),
		authHandler:         NewAuthHandler(authService, log),
		businessHandler:     NewBusinessHandler(businessService, log),
		callHandler:         NewCallHandler(callService, log),
		analyticsHandler:    NewAnalyticsHandler(analyticsService, log),
		interactionHandler:  NewInteractionHandler(interactionService, log),
	}

	r.setupRoutes()
	return r
}

func (r *Router) setupRoutes() {
	// Apply global middleware
	r.router.Use(r.errorMiddleware.Recovery)
	r.router.Use(r.loggingMiddleware.Log)
	r.router.Use(r.corsMiddleware.Handle)

	// Health check endpoints (no auth required)
	r.router.HandleFunc("/health", r.healthCheck).Methods("GET")
	r.router.HandleFunc("/ready", r.readyCheck).Methods("GET")

	// API v1 routes
	api := r.router.PathPrefix("/api/v1").Subrouter()

	// Public routes (no auth required)
	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", r.authHandler.Register).Methods("POST")
	auth.HandleFunc("/login", r.authHandler.Login).Methods("POST")
	auth.HandleFunc("/refresh", r.authHandler.RefreshToken).Methods("POST")

	// Webhook route (no auth - validated by signature)
	api.HandleFunc("/webhooks/vapi", r.callHandler.HandleWebhook).Methods("POST")

	// Protected routes (require authentication)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(r.authMiddleware.Authenticate)

	// Auth routes
	protected.HandleFunc("/auth/logout", r.authHandler.Logout).Methods("POST")

	// Business routes
	protected.HandleFunc("/businesses/me", r.businessHandler.GetBusiness).Methods("GET")
	protected.HandleFunc("/businesses/me", r.businessHandler.UpdateBusiness).Methods("PUT")

	// Call routes
	protected.HandleFunc("/calls", r.callHandler.InitiateCall).Methods("POST")
	protected.HandleFunc("/calls", r.callHandler.ListCalls).Methods("GET")
	protected.HandleFunc("/calls/{id}", r.callHandler.GetCall).Methods("GET")
	protected.HandleFunc("/calls/{id}/transcript", r.callHandler.GetTranscript).Methods("GET")
	protected.HandleFunc("/calls/{id}/interactions", r.interactionHandler.GetCallInteractions).Methods("GET")

	// Interaction routes
	protected.HandleFunc("/interactions", r.interactionHandler.ListInteractions).Methods("GET")

	// Appointment routes
	protected.HandleFunc("/appointments", r.interactionHandler.ListAppointments).Methods("GET")
	protected.HandleFunc("/appointments/{id}", r.interactionHandler.UpdateAppointmentStatus).Methods("PATCH")

	// Analytics routes
	protected.HandleFunc("/analytics/overview", r.analyticsHandler.GetOverview).Methods("GET")
	protected.HandleFunc("/analytics/calls", r.analyticsHandler.GetCallVolume).Methods("GET")
}

func (r *Router) healthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func (r *Router) readyCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ready"}`))
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
