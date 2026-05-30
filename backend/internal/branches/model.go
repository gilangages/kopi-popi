package branches

import "time"

// Branch merepresentasikan struktur data tabel branches di MySQL
type Branch struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateBranchRequest merepresentasikan payload untuk membuat cabang baru
type CreateBranchRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
}

// UpdateBranchRequest merepresentasikan payload untuk memperbarui cabang
type UpdateBranchRequest struct {
	Name     *string `json:"name,omitempty"`
	Address  *string `json:"address,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}
