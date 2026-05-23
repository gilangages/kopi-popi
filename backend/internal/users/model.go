package users

import "time"

// Role merepresentasikan tabel roles
type Role struct {
	ID   int    `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

// Branch merepresentasikan tabel branches
type Branch struct {
	ID   int    `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

// User merepresentasikan struktur data dari database (mirip dengan auth.User tapi ini ranah Users domain)
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

// UserResponse adalah response umum untuk data user
type UserResponse struct {
	ID             string  `json:"id"`
	RoleID         int     `json:"role_id"`
	BranchID       *int    `json:"branch_id,omitempty"`
	Name           string  `json:"name"`
	Email          string  `json:"email"`
	Phone          *string `json:"phone,omitempty"`
	ProfilePicture string  `json:"profile_picture"`
	IsActive       bool    `json:"is_active"`
}

// UpdateProfileRequest DTO untuk /users/me
type UpdateProfileRequest struct {
	Name           *string `json:"name,omitempty"`
	Phone          *string `json:"phone,omitempty"`
	ProfilePicture *string `json:"profile_picture,omitempty"`
}

// UpdatePasswordRequest DTO untuk /users/me/password
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}

// RequestEmailOTPRequest DTO
type RequestEmailOTPRequest struct {
	NewEmail string `json:"new_email" binding:"required,email"`
}

// VerifyEmailOTPRequest DTO
type VerifyEmailOTPRequest struct {
	OTP      string `json:"otp" binding:"required"`
	NewEmail string `json:"new_email" binding:"required,email"`
}

// CreateManagerRequest DTO (Admin Only)
type CreateManagerRequest struct {
	BranchID int    `json:"branch_id" binding:"required"` // Asumsi int sesuai migration database (000003_create_users_table)
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// CreateCashierRequest DTO (Manager Only)
// Manager tidak mengirim BranchID, kita inject dari backend!
type CreateCashierRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}
