package users

import (
	"context"
	"errors"
	"log"

	"github.com/gilangages/kopi-popi/pkg/hash"
	"github.com/google/uuid"
)

type Service interface {
	GetMyProfile(ctx context.Context, userID string) (*UserResponse, error)
	UpdateMyProfile(ctx context.Context, userID string, req UpdateProfileRequest) error
	DeleteProfilePicture(ctx context.Context, userID string) error
	UpdateMyPassword(ctx context.Context, userID string, req UpdatePasswordRequest) error
	RequestEmailOTP(ctx context.Context, userID string, req RequestEmailOTPRequest) error
	VerifyEmailOTP(ctx context.Context, userID string, req VerifyEmailOTPRequest) error
	
	CreateManager(ctx context.Context, req CreateManagerRequest) error
	CreateCashier(ctx context.Context, managerBranchID int, req CreateCashierRequest) error
	GetEmployees(ctx context.Context, role string, branchID *int) ([]UserResponse, error)
	DisableEmployee(ctx context.Context, employeeID string) error
}

type userService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &userService{repo}
}

func (s *userService) GetMyProfile(ctx context.Context, userID string) (*UserResponse, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return &UserResponse{
		ID:             user.ID,
		RoleID:         user.RoleID,
		BranchID:       user.BranchID,
		Name:           user.Name,
		Email:          user.Email,
		Phone:          user.Phone,
		ProfilePicture: user.ProfilePicture,
		IsActive:       user.IsActive,
	}, nil
}

func (s *userService) UpdateMyProfile(ctx context.Context, userID string, req UpdateProfileRequest) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.ProfilePicture != nil {
		user.ProfilePicture = *req.ProfilePicture
	}

	return s.repo.Update(ctx, user)
}

func (s *userService) DeleteProfilePicture(ctx context.Context, userID string) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Hapus foto profil dengan mengatur nilainya menjadi string kosong
	user.ProfilePicture = ""
	return s.repo.Update(ctx, user)
}

func (s *userService) UpdateMyPassword(ctx context.Context, userID string, req UpdatePasswordRequest) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Verifikasi current password
	if !hash.CheckHash(req.CurrentPassword, user.PasswordHash) {
		return errors.New("current password is incorrect")
	}

	hashedPassword, err := hash.MakeHash(req.NewPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = hashedPassword
	return s.repo.Update(ctx, user)
}

func (s *userService) RequestEmailOTP(ctx context.Context, userID string, req RequestEmailOTPRequest) error {
	// Di sistem nyata, kita akan generate OTP 6 digit, simpan ke database OTP table / Redis,
	// lalu mengirimnya ke req.NewEmail via Email Service (SMTP).
	// Untuk demo ini, kita asumsikan OTP 123456 selalu dikirim.
	log.Printf("MENGIRIM OTP (123456) KE EMAIL BARU: %s", req.NewEmail)
	return nil
}

func (s *userService) VerifyEmailOTP(ctx context.Context, userID string, req VerifyEmailOTPRequest) error {
	// Asumsi validasi OTP statis
	if req.OTP != "123456" {
		return errors.New("invalid OTP")
	}

	// Pastikan email belum dipakai oleh orang lain
	existingUser, err := s.repo.FindByEmail(ctx, req.NewEmail)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("email is already used by another account")
	}

	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	user.Email = req.NewEmail
	return s.repo.Update(ctx, user)
}

func (s *userService) CreateManager(ctx context.Context, req CreateManagerRequest) error {
	// 1. Cek apakah cabang ada
	exists, err := s.repo.CheckBranchExists(ctx, req.BranchID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("branch does not exist")
	}

	// 2. Cek aturan 1 Cabang = 1 Manager
	hasManager, err := s.repo.CheckBranchManagerExists(ctx, req.BranchID)
	if err != nil {
		return err
	}
	if hasManager {
		return errors.New("this branch already has a manager")
	}

	// 3. Pastikan email belum terpakai
	existingUser, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("email is already registered")
	}

	// 4. Buat User Manager
	hashedPassword, err := hash.MakeHash(req.Password)
	if err != nil {
		return err
	}

	user := &User{
		ID:           uuid.New().String(),
		RoleID:       2, // 2 = Manager
		BranchID:     &req.BranchID,
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		IsActive:     true,
	}

	return s.repo.CreateUser(ctx, user)
}

func (s *userService) CreateCashier(ctx context.Context, managerBranchID int, req CreateCashierRequest) error {
	// 1. Pastikan email belum terpakai
	existingUser, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("email is already registered")
	}

	// 2. Buat User Kasir dan otomatis assign ke cabang manager
	hashedPassword, err := hash.MakeHash(req.Password)
	if err != nil {
		return err
	}

	user := &User{
		ID:           uuid.New().String(),
		RoleID:       3, // 3 = Cashier
		BranchID:     &managerBranchID,
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		IsActive:     true,
	}

	return s.repo.CreateUser(ctx, user)
}

func (s *userService) GetEmployees(ctx context.Context, role string, branchID *int) ([]UserResponse, error) {
	var users []User
	var err error

	if role == "Admin" {
		users, err = s.repo.FindAllEmployees(ctx)
	} else if role == "Manager" {
		if branchID == nil {
			return nil, errors.New("manager does not have a branch assigned")
		}
		users, err = s.repo.FindEmployeesByBranch(ctx, *branchID)
	} else {
		return nil, errors.New("forbidden: insufficient privileges")
	}

	if err != nil {
		return nil, err
	}

	responses := []UserResponse{} // Inisialisasi slice kosong agar hasilnya [] bukan null
	for _, u := range users {
		responses = append(responses, UserResponse{
			ID:             u.ID,
			RoleID:         u.RoleID,
			BranchID:       u.BranchID,
			Name:           u.Name,
			Email:          u.Email,
			Phone:          u.Phone,
			ProfilePicture: u.ProfilePicture,
			IsActive:       u.IsActive,
		})
	}
	return responses, nil
}

func (s *userService) DisableEmployee(ctx context.Context, employeeID string) error {
	user, err := s.repo.FindByID(ctx, employeeID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("employee not found")
	}

	// Toggle IsActive (jika tadinya true jadi false, jika false jadi true)
	user.IsActive = !user.IsActive
	return s.repo.Update(ctx, user)
}
