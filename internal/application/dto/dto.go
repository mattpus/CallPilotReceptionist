package dto

// Authentication DTOs

type RegisterRequest struct {
	BusinessName string `json:"business_name"`
	BusinessType string `json:"business_type"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Password     string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// User DTOs

type UserResponse struct {
	ID         string `json:"id"`
	BusinessID string `json:"business_id"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	CreatedAt  string `json:"created_at"`
}

// Business DTOs

type BusinessResponse struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Phone     string                 `json:"phone"`
	Settings  map[string]interface{} `json:"settings"`
	CreatedAt string                 `json:"created_at"`
	UpdatedAt string                 `json:"updated_at"`
}

type UpdateBusinessRequest struct {
	Name     string                 `json:"name,omitempty"`
	Type     string                 `json:"type,omitempty"`
	Phone    string                 `json:"phone,omitempty"`
	Settings map[string]interface{} `json:"settings,omitempty"`
}

// Call DTOs

type InitiateCallRequest struct {
	PhoneNumber string                 `json:"phone_number"`
	AssistantID string                 `json:"assistant_id,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type CallResponse struct {
	ID             string                 `json:"id"`
	BusinessID     string                 `json:"business_id"`
	ProviderCallID string                 `json:"provider_call_id,omitempty"`
	CallerPhone    string                 `json:"caller_phone"`
	Duration       int                    `json:"duration"`
	Status         string                 `json:"status"`
	Cost           float64                `json:"cost"`
	StartedAt      *string                `json:"started_at,omitempty"`
	EndedAt        *string                `json:"ended_at,omitempty"`
	CreatedAt      string                 `json:"created_at"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

type ListCallsRequest struct {
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Status string `json:"status,omitempty"`
}

type ListCallsResponse struct {
	Calls      []CallResponse `json:"calls"`
	Total      int            `json:"total"`
	Limit      int            `json:"limit"`
	Offset     int            `json:"offset"`
}

// Interaction DTOs

type InteractionResponse struct {
	ID        string                 `json:"id"`
	CallID    string                 `json:"call_id"`
	Type      string                 `json:"type"`
	Content   map[string]interface{} `json:"content"`
	Timestamp string                 `json:"timestamp"`
	CreatedAt string                 `json:"created_at"`
}

type ListInteractionsResponse struct {
	Interactions []InteractionResponse `json:"interactions"`
	Total        int                   `json:"total"`
	Limit        int                   `json:"limit"`
	Offset       int                   `json:"offset"`
}

// Transcript DTOs

type TranscriptMessageResponse struct {
	Role      string `json:"role"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

type TranscriptResponse struct {
	CallID   string                      `json:"call_id"`
	Messages []TranscriptMessageResponse `json:"messages"`
}

// Appointment DTOs

type AppointmentResponse struct {
	ID             string  `json:"id"`
	CallID         string  `json:"call_id"`
	BusinessID     string  `json:"business_id"`
	CustomerName   string  `json:"customer_name,omitempty"`
	CustomerPhone  string  `json:"customer_phone"`
	RequestedDate  *string `json:"requested_date,omitempty"`
	RequestedTime  string  `json:"requested_time,omitempty"`
	ServiceType    string  `json:"service_type,omitempty"`
	Notes          string  `json:"notes,omitempty"`
	Status         string  `json:"status"`
	ExtractedAt    string  `json:"extracted_at"`
	ConfirmedAt    *string `json:"confirmed_at,omitempty"`
	CreatedAt      string  `json:"created_at"`
}

type UpdateAppointmentRequest struct {
	Status string `json:"status"`
}

// Analytics DTOs

type AnalyticsOverviewResponse struct {
	TotalCalls      int     `json:"total_calls"`
	CompletedCalls  int     `json:"completed_calls"`
	FailedCalls     int     `json:"failed_calls"`
	TotalDuration   int     `json:"total_duration"`   // seconds
	AverageDuration float64 `json:"average_duration"` // seconds
	TotalCost       float64 `json:"total_cost"`
	PendingAppointments int `json:"pending_appointments"`
}

type CallVolumeData struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type CallVolumeResponse struct {
	Data []CallVolumeData `json:"data"`
}

// Error response

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Success response

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
