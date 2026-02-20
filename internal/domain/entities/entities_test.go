package entities

import (
	"testing"
	"time"
)

func TestNewBusiness(t *testing.T) {
	tests := []struct {
		name         string
		businessName string
		businessType string
		phone        string
		settings     map[string]interface{}
		wantErr      bool
	}{
		{
			name:         "valid business",
			businessName: "Test Business",
			businessType: "dentist",
			phone:        "+1234567890",
			settings:     nil,
			wantErr:      false,
		},
		{
			name:         "missing name",
			businessName: "",
			businessType: "dentist",
			phone:        "+1234567890",
			settings:     nil,
			wantErr:      true,
		},
		{
			name:         "missing phone",
			businessName: "Test Business",
			businessType: "dentist",
			phone:        "",
			settings:     nil,
			wantErr:      true,
		},
		{
			name:         "with settings",
			businessName: "Test Business",
			businessType: "dentist",
			phone:        "+1234567890",
			settings:     map[string]interface{}{"key": "value"},
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			business, err := NewBusiness(tt.businessName, tt.businessType, tt.phone, tt.settings)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewBusiness() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if business.Name != tt.businessName {
					t.Errorf("NewBusiness() name = %v, want %v", business.Name, tt.businessName)
				}
				if business.Phone != tt.phone {
					t.Errorf("NewBusiness() phone = %v, want %v", business.Phone, tt.phone)
				}
			}
		})
	}
}

func TestBusiness_Update(t *testing.T) {
	business, _ := NewBusiness("Original", "dentist", "+1234567890", nil)
	originalUpdatedAt := business.UpdatedAt

	time.Sleep(10 * time.Millisecond)

	err := business.Update("Updated", "clinic", "+0987654321", map[string]interface{}{"new": "setting"})

	if err != nil {
		t.Errorf("Update() unexpected error: %v", err)
	}

	if business.Name != "Updated" {
		t.Errorf("Update() name = %v, want Updated", business.Name)
	}

	if business.UpdatedAt.Equal(originalUpdatedAt) {
		t.Error("Update() UpdatedAt should be updated")
	}
}

func TestNewCall(t *testing.T) {
	tests := []struct {
		name        string
		businessID  string
		callerPhone string
		wantErr     bool
	}{
		{
			name:        "valid call",
			businessID:  "business-123",
			callerPhone: "+1234567890",
			wantErr:     false,
		},
		{
			name:        "missing business ID",
			businessID:  "",
			callerPhone: "+1234567890",
			wantErr:     true,
		},
		{
			name:        "missing caller phone",
			businessID:  "business-123",
			callerPhone: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			call, err := NewCall(tt.businessID, tt.callerPhone)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewCall() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if call.BusinessID != tt.businessID {
					t.Errorf("NewCall() businessID = %v, want %v", call.BusinessID, tt.businessID)
				}
				if call.Status != CallStatusInitiated {
					t.Errorf("NewCall() status = %v, want %v", call.Status, CallStatusInitiated)
				}
			}
		})
	}
}

func TestCall_UpdateStatus(t *testing.T) {
	call, _ := NewCall("business-123", "+1234567890")

	tests := []struct {
		name       string
		status     CallStatus
		wantErr    bool
		checkField func(*testing.T, *Call)
	}{
		{
			name:    "update to ringing",
			status:  CallStatusRinging,
			wantErr: false,
			checkField: func(t *testing.T, c *Call) {
				if c.Status != CallStatusRinging {
					t.Error("Status not updated to ringing")
				}
			},
		},
		{
			name:    "update to in_progress",
			status:  CallStatusInProgress,
			wantErr: false,
			checkField: func(t *testing.T, c *Call) {
				if c.StartedAt == nil {
					t.Error("StartedAt should be set")
				}
			},
		},
		{
			name:    "update to completed",
			status:  CallStatusCompleted,
			wantErr: false,
			checkField: func(t *testing.T, c *Call) {
				if c.EndedAt == nil {
					t.Error("EndedAt should be set")
				}
			},
		},
		{
			name:    "invalid status",
			status:  CallStatus("invalid"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := call.UpdateStatus(tt.status)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.checkField != nil {
				tt.checkField(t, call)
			}
		})
	}
}

func TestCall_IsCompleted(t *testing.T) {
	tests := []struct {
		name   string
		status CallStatus
		want   bool
	}{
		{"completed", CallStatusCompleted, true},
		{"failed", CallStatusFailed, true},
		{"no answer", CallStatusNoAnswer, true},
		{"busy", CallStatusBusy, true},
		{"initiated", CallStatusInitiated, false},
		{"in progress", CallStatusInProgress, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			call, _ := NewCall("business-123", "+1234567890")
			call.UpdateStatus(tt.status)

			if got := call.IsCompleted(); got != tt.want {
				t.Errorf("IsCompleted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewUser(t *testing.T) {
	tests := []struct {
		name         string
		businessID   string
		email        string
		passwordHash string
		role         UserRole
		wantErr      bool
	}{
		{
			name:         "valid user",
			businessID:   "business-123",
			email:        "test@example.com",
			passwordHash: "hashed_password",
			role:         UserRoleOwner,
			wantErr:      false,
		},
		{
			name:         "missing business ID",
			businessID:   "",
			email:        "test@example.com",
			passwordHash: "hashed_password",
			role:         UserRoleOwner,
			wantErr:      true,
		},
		{
			name:         "invalid role",
			businessID:   "business-123",
			email:        "test@example.com",
			passwordHash: "hashed_password",
			role:         UserRole("invalid"),
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.businessID, tt.email, tt.passwordHash, tt.role)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if user.BusinessID != tt.businessID {
					t.Errorf("NewUser() businessID = %v, want %v", user.BusinessID, tt.businessID)
				}
				if user.Role != tt.role {
					t.Errorf("NewUser() role = %v, want %v", user.Role, tt.role)
				}
			}
		})
	}
}

func TestUser_Permissions(t *testing.T) {
	tests := []struct {
		name           string
		role           UserRole
		isOwner        bool
		isAdmin        bool
		canManage      bool
	}{
		{"owner", UserRoleOwner, true, true, true},
		{"admin", UserRoleAdmin, false, true, true},
		{"employee", UserRoleEmployee, false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, _ := NewUser("business-123", "test@example.com", "hash", tt.role)

			if user.IsOwner() != tt.isOwner {
				t.Errorf("IsOwner() = %v, want %v", user.IsOwner(), tt.isOwner)
			}
			if user.IsAdmin() != tt.isAdmin {
				t.Errorf("IsAdmin() = %v, want %v", user.IsAdmin(), tt.isAdmin)
			}
			if user.CanManageBusiness() != tt.canManage {
				t.Errorf("CanManageBusiness() = %v, want %v", user.CanManageBusiness(), tt.canManage)
			}
		})
	}
}

func TestNewAppointmentRequest(t *testing.T) {
	requestedDate := time.Now().Add(24 * time.Hour)

	tests := []struct {
		name          string
		callID        string
		businessID    string
		customerPhone string
		wantErr       bool
	}{
		{
			name:          "valid appointment",
			callID:        "call-123",
			businessID:    "business-123",
			customerPhone: "+1234567890",
			wantErr:       false,
		},
		{
			name:          "missing call ID",
			callID:        "",
			businessID:    "business-123",
			customerPhone: "+1234567890",
			wantErr:       true,
		},
		{
			name:          "missing customer phone",
			callID:        "call-123",
			businessID:    "business-123",
			customerPhone: "",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apt, err := NewAppointmentRequest(
				tt.callID,
				tt.businessID,
				"John Doe",
				tt.customerPhone,
				&requestedDate,
				"10:00 AM",
				"cleaning",
				"notes",
			)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewAppointmentRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if apt.Status != AppointmentStatusPending {
					t.Errorf("NewAppointmentRequest() status = %v, want %v", apt.Status, AppointmentStatusPending)
				}
			}
		})
	}
}

func TestAppointmentRequest_StatusTransitions(t *testing.T) {
	apt, _ := NewAppointmentRequest(
		"call-123",
		"business-123",
		"John Doe",
		"+1234567890",
		nil,
		"10:00 AM",
		"cleaning",
		"",
	)

	// Test confirm
	if err := apt.Confirm(); err != nil {
		t.Errorf("Confirm() unexpected error: %v", err)
	}
	if apt.Status != AppointmentStatusConfirmed {
		t.Error("Status should be confirmed")
	}
	if apt.ConfirmedAt == nil {
		t.Error("ConfirmedAt should be set")
	}

	// Test complete
	if err := apt.Complete(); err != nil {
		t.Errorf("Complete() unexpected error: %v", err)
	}
	if apt.Status != AppointmentStatusCompleted {
		t.Error("Status should be completed")
	}

	// Test invalid transition (cancel completed)
	if err := apt.Cancel(); err == nil {
		t.Error("Cancel() should fail on completed appointment")
	}
}
