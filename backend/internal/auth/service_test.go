package auth

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/gilangages/kopi-popi/pkg/hash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- MOCK REPOSITORY ---
type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) CreateUser(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) != nil {
		return args.Get(0).(*User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRepository) CreatePasswordReset(ctx context.Context, pwReset *PasswordReset) error {
	args := m.Called(ctx, pwReset)
	return args.Error(0)
}

func (m *mockRepository) FindPasswordResetByToken(ctx context.Context, token string) (*PasswordReset, error) {
	args := m.Called(ctx, token)
	if args.Get(0) != nil {
		return args.Get(0).(*PasswordReset), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRepository) UpdatePassword(ctx context.Context, email string, hashedPassword string) error {
	args := m.Called(ctx, email, hashedPassword)
	return args.Error(0)
}

func (m *mockRepository) DeletePasswordReset(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

// Inisialisasi env sebelum testing
func init() {
	os.Setenv("JWT_SECRET", "secret_untuk_testing")
}

// --- UNIT TESTS ---

func TestRegister_Success(t *testing.T) {
	mockRepo := new(mockRepository)
	service := NewService(mockRepo)

	req := RegisterRequest{
		Name:            "Gilang",
		Email:           "gilang@example.com",
		Password:        "rahasia123",
		ConfirmPassword: "rahasia123",
		Phone:           "0812345",
	}

	// Skenario: FindByEmail mengembalikan nil (email belum ada)
	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, nil)
	// Skenario: CreateUser mengembalikan sukses (nil)
	mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*auth.User")).Return(nil)

	user, err := service.Register(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.Name, user.Name)
	assert.Equal(t, req.Email, user.Email)
	mockRepo.AssertExpectations(t)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	mockRepo := new(mockRepository)
	service := NewService(mockRepo)

	req := RegisterRequest{
		Name:     "Gilang",
		Email:    "gilang@example.com",
		Password: "123",
	}

	existingUser := &User{Email: "gilang@example.com"}
	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(existingUser, nil)

	user, err := service.Register(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "email is already registered", err.Error())
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	mockRepo := new(mockRepository)
	service := NewService(mockRepo)

	password := "rahasia123"
	hashedPassword, _ := hash.MakeHash(password)

	user := &User{
		ID:           "uuid-123",
		Name:         "Gilang",
		Email:        "gilang@example.com",
		PasswordHash: hashedPassword,
		RoleID:       4, // Customer
	}

	req := LoginRequest{
		Email:      "gilang@example.com",
		Password:   "rahasia123",
		RememberMe: true,
	}

	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(user, nil)

	token, returnedUser, err := service.Login(context.Background(), req)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, user.ID, returnedUser.ID)
	mockRepo.AssertExpectations(t)
}

func TestLogin_InvalidEmailOrPassword(t *testing.T) {
	mockRepo := new(mockRepository)
	service := NewService(mockRepo)

	// Skenario 1: Email tidak ada
	mockRepo.On("FindByEmail", mock.Anything, "salah@example.com").Return(nil, nil).Once()

	_, _, err := service.Login(context.Background(), LoginRequest{Email: "salah@example.com", Password: "123"})
	assert.Error(t, err)
	assert.Equal(t, "invalid email or password", err.Error())

	// Skenario 2: Password salah
	password := "rahasia123"
	hashedPassword, _ := hash.MakeHash(password)
	user := &User{Email: "gilang@example.com", PasswordHash: hashedPassword}

	mockRepo.On("FindByEmail", mock.Anything, "gilang@example.com").Return(user, nil).Once()

	_, _, err2 := service.Login(context.Background(), LoginRequest{Email: "gilang@example.com", Password: "salah_password"})
	assert.Error(t, err2)
	assert.Equal(t, "invalid email or password", err2.Error())

	mockRepo.AssertExpectations(t)
}

func TestForgotPassword_Success(t *testing.T) {
	mockRepo := new(mockRepository)
	service := NewService(mockRepo)

	req := ForgotPasswordRequest{Email: "gilang@example.com"}
	user := &User{Email: "gilang@example.com"}

	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(user, nil)
	mockRepo.On("DeletePasswordReset", mock.Anything, req.Email).Return(nil)
	mockRepo.On("CreatePasswordReset", mock.Anything, mock.AnythingOfType("*auth.PasswordReset")).Return(nil)

	err := service.ForgotPassword(context.Background(), req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestResetPassword_Success(t *testing.T) {
	mockRepo := new(mockRepository)
	service := NewService(mockRepo)

	req := ResetPasswordRequest{
		Token:           "valid-token",
		NewPassword:     "passwordbaru",
		ConfirmPassword: "passwordbaru",
	}

	pwReset := &PasswordReset{
		Email:     "gilang@example.com",
		Token:     "valid-token",
		ExpiresAt: time.Now().Add(time.Hour * 1), // Belum expired
	}

	mockRepo.On("FindPasswordResetByToken", mock.Anything, req.Token).Return(pwReset, nil)
	mockRepo.On("UpdatePassword", mock.Anything, pwReset.Email, mock.AnythingOfType("string")).Return(nil)
	mockRepo.On("DeletePasswordReset", mock.Anything, pwReset.Email).Return(nil)

	err := service.ResetPassword(context.Background(), req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestResetPassword_ExpiredToken(t *testing.T) {
	mockRepo := new(mockRepository)
	service := NewService(mockRepo)

	req := ResetPasswordRequest{Token: "expired-token", NewPassword: "123", ConfirmPassword: "123"}

	pwReset := &PasswordReset{
		Email:     "gilang@example.com",
		Token:     "expired-token",
		ExpiresAt: time.Now().Add(-time.Hour * 1), // Sudah expired 1 jam lalu
	}

	mockRepo.On("FindPasswordResetByToken", mock.Anything, req.Token).Return(pwReset, nil)
	mockRepo.On("DeletePasswordReset", mock.Anything, pwReset.Email).Return(nil)

	err := service.ResetPassword(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "invalid or expired token", err.Error())
	mockRepo.AssertExpectations(t)
}
