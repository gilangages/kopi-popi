package auth

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/gilangages/kopi-popi/pkg/hash"
	"github.com/gilangages/kopi-popi/pkg/jwt"
	"github.com/google/uuid"
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (*User, error)
	Login(ctx context.Context, req LoginRequest) (string, *User, error)
	ForgotPassword(ctx context.Context, req ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req ResetPasswordRequest) error
}

type authService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &authService{repo}
}

func (s *authService) Register(ctx context.Context, req RegisterRequest) (*User, error) {
	// 1. Cek ketersediaan email
	existingUser, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email is already registered")
	}

	// 2. Hash password
	hashedPassword, err := hash.MakeHash(req.Password)
	if err != nil {
		return nil, err
	}

	// 3. Buat User ID baru
	newID := uuid.New().String()

	// 4. Siapkan object user
	var phonePtr *string
	if req.Phone != "" {
		phonePtr = &req.Phone
	}

	user := &User{
		ID:           newID,
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Phone:        phonePtr,
	}

	// 5. Simpan ke database
	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, req LoginRequest) (string, *User, error) {
	// 1. Cari user berdasarkan email
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, errors.New("invalid email or password")
	}

	// 2. Cocokkan hash password
	isValid := hash.CheckHash(req.Password, user.PasswordHash)
	if !isValid {
		return "", nil, errors.New("invalid email or password")
	}

	// 3. Generate JWT
	// Asumsi default role name
	roleName := "Customer"
	if user.RoleID == 1 { // Asumsi ID 1 adalah Admin
		roleName = "Admin"
	} else if user.RoleID == 2 {
		roleName = "Manager"
	} else if user.RoleID == 3 {
		roleName = "Cashier"
	}

	token, err := jwt.GenerateToken(user.ID, user.Name, roleName, user.BranchID, req.RememberMe)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *authService) ForgotPassword(ctx context.Context, req ForgotPasswordRequest) error {
	// 1. Cek apakah email terdaftar
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if user == nil {
		// Security best practice: Jangan beri tahu email tidak ada,
		// pura-pura sukses agar tidak bisa ditebak (enumeration)
		return nil
	}

	// 2. Generate Reset Token (UUID biasa sudah cukup aman)
	resetToken := uuid.New().String()

	// 3. Simpan ke Database
	pwReset := &PasswordReset{
		Email:     req.Email,
		Token:     resetToken,
		ExpiresAt: time.Now().Add(time.Hour * 1), // Berlaku 1 jam
	}

	// Hapus token lama jika ada agar tidak menumpuk
	_ = s.repo.DeletePasswordReset(ctx, req.Email)

	err = s.repo.CreatePasswordReset(ctx, pwReset)
	if err != nil {
		return err
	}

	// 4. Simulasi pengiriman email
	// Di sistem nyata, panggil fungsi pengirim email (misal via SMTP / SendGrid)
	log.Printf("\n======================================================\n")
	log.Printf("MENGIRIM EMAIL KE: %s\n", req.Email)
	log.Printf("LINK RESET PASSWORD: http://localhost:3000/reset-password?token=%s\n", resetToken)
	log.Printf("======================================================\n\n")
	
	return nil
}

func (s *authService) ResetPassword(ctx context.Context, req ResetPasswordRequest) error {
	// 1. Cari token di database
	pwReset, err := s.repo.FindPasswordResetByToken(ctx, req.Token)
	if err != nil {
		return err
	}
	if pwReset == nil {
		return errors.New("invalid or expired token")
	}

	// 2. Cek kedaluwarsa
	if time.Now().After(pwReset.ExpiresAt) {
		_ = s.repo.DeletePasswordReset(ctx, pwReset.Email) // Bersihkan token
		return errors.New("invalid or expired token")
	}

	// 3. Hash password baru
	hashedPassword, err := hash.MakeHash(req.NewPassword)
	if err != nil {
		return err
	}

	// 4. Update password di tabel users
	err = s.repo.UpdatePassword(ctx, pwReset.Email, hashedPassword)
	if err != nil {
		return err
	}

	// 5. Hapus token agar tidak bisa dipakai 2x
	_ = s.repo.DeletePasswordReset(ctx, pwReset.Email)

	return nil
}
