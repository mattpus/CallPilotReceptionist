package entities

import (
	"time"

	"github.com/CallPilotReceptionist/internal/domain/errors"
)

type UserRole string

const (
	UserRoleOwner    UserRole = "owner"
	UserRoleAdmin    UserRole = "admin"
	UserRoleEmployee UserRole = "employee"
)

type User struct {
	ID           string    `json:"id"`
	BusinessID   string    `json:"business_id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

func NewUser(businessID, email, passwordHash string, role UserRole) (*User, error) {
	if businessID == "" {
		return nil, errors.NewValidationError("business_id is required")
	}
	if email == "" {
		return nil, errors.NewValidationError("email is required")
	}
	if passwordHash == "" {
		return nil, errors.NewValidationError("password is required")
	}

	validRoles := map[UserRole]bool{
		UserRoleOwner:    true,
		UserRoleAdmin:    true,
		UserRoleEmployee: true,
	}

	if !validRoles[role] {
		return nil, errors.NewValidationError("invalid user role")
	}

	return &User{
		BusinessID:   businessID,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
		CreatedAt:    time.Now(),
	}, nil
}

func (u *User) IsOwner() bool {
	return u.Role == UserRoleOwner
}

func (u *User) IsAdmin() bool {
	return u.Role == UserRoleAdmin || u.Role == UserRoleOwner
}

func (u *User) CanManageBusiness() bool {
	return u.IsAdmin()
}

func (u *User) Validate() error {
	if u.BusinessID == "" {
		return errors.NewValidationError("business_id is required")
	}
	if u.Email == "" {
		return errors.NewValidationError("email is required")
	}
	if u.PasswordHash == "" {
		return errors.NewValidationError("password is required")
	}
	return nil
}
