package auth

import "time"

// Role merepresentasikan tabel roles di database
type Role struct {
	ID   int    `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

// Branch merepresentasikan tabel branches di database
type Branch struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// User merepresentasikan tabel users di database
type User struct {
	ID             string    `json:"id" gorm:"primaryKey"`
	RoleID         int       `json:"role_id"`
	Role           Role      `json:"role" gorm:"foreignKey:RoleID"`
	BranchID       *int      `json:"branch_id,omitempty"`
	Branch         *Branch   `json:"branch,omitempty" gorm:"foreignKey:BranchID"`
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

// LoginRequest adalah DTO untuk menerima input saat login
type LoginRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	RememberMe bool   `json:"remember_me"`
}

// PasswordReset merepresentasikan tabel password_resets di database
type PasswordReset struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// ForgotPasswordRequest adalah DTO untuk meminta reset token
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest adalah DTO untuk mengubah password dengan token
type ResetPasswordRequest struct {
	Token           string `json:"token" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}
