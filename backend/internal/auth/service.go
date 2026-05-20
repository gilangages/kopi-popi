package auth

import (
	"context"
	"errors"

	"github.com/gilangages/kopi-popi/pkg/hash"
	"github.com/gilangages/kopi-popi/pkg/jwt"
	"github.com/google/uuid"
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (*User, error)
	Login(ctx context.Context, req LoginRequest) (string, *User, error)
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

	token, err := jwt.GenerateToken(user.ID, user.Name, roleName)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}
