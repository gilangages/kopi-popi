package auth

import "time"

// User merepresentasikan tabel users di database
type User struct {
	ID             string    `json:"id"`
	RoleID         int       `json:"role_id"`
	BranchID       *int      `json:"branch_id,omitempty"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	PasswordHash   string    `json:"-"`
	Phone          *string   `json:"phone,omitempty"`
	ProfilePicture string    `json:"profile_picture"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// RegisterRequest adalah DTO untuk menerima input dari client saat register
type RegisterRequest struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
	Phone           string `json:"phone"`
}

// LoginRequest adalah DTO untuk menerima input dari client saat login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
