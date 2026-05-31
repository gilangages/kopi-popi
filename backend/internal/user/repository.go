package user

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	CheckBranchManagerExists(ctx context.Context, branchID int) (bool, error)
	CheckBranchExists(ctx context.Context, branchID int) (bool, error)
	CreateUser(ctx context.Context, user *User) error
	FindAllEmployees(ctx context.Context) ([]User, error)
	FindEmployeesByBranch(ctx context.Context, branchID int) ([]User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &userRepository{db}
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) CheckBranchExists(ctx context.Context, branchID int) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("branches").Where("id = ?", branchID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userRepository) CheckBranchManagerExists(ctx context.Context, branchID int) (bool, error) {
	var count int64
	// Role 2 adalah Manager
	err := r.db.WithContext(ctx).Model(&User{}).Where("branch_id = ? AND role_id = ?", branchID, 2).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindAllEmployees(ctx context.Context) ([]User, error) {
	var users []User
	// Mengambil semua user (termasuk Admin & Customer)
	err := r.db.WithContext(ctx).Find(&users).Error
	return users, err
}

func (r *userRepository) FindEmployeesByBranch(ctx context.Context, branchID int) ([]User, error) {
	var users []User
	// Hanya mencari kasir (Role 3) di cabangnya sendiri. Manager tak perlu lihat profilnya sendiri dari list.
	// Jika mau menampilkan manager juga: Where("branch_id = ? AND role_id IN (2, 3)", branchID)
	err := r.db.WithContext(ctx).Where("branch_id = ? AND role_id = 3", branchID).Find(&users).Error
	return users, err
}
