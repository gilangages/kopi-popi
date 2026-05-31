package user

import (
	"context"
	"testing"

	"github.com/gilangages/kopi-popi/pkg/hash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- MOCK REPOSITORY ---
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) FindByID(ctx context.Context, id string) (*User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) != nil {
		return args.Get(0).(*User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) CheckBranchExists(ctx context.Context, branchID int) (bool, error) {
	args := m.Called(ctx, branchID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) CheckBranchManagerExists(ctx context.Context, branchID int) (bool, error) {
	args := m.Called(ctx, branchID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) CreateUser(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) FindAllEmployees(ctx context.Context) ([]User, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) FindEmployeesByBranch(ctx context.Context, branchID int) ([]User, error) {
	args := m.Called(ctx, branchID)
	if args.Get(0) != nil {
		return args.Get(0).([]User), args.Error(1)
	}
	return nil, args.Error(1)
}

// --- PENGUJIAN PROFIL ---

func TestGetMyProfile_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	mockUser := &User{
		ID:       "user-123",
		RoleID:   1,
		Name:     "Gilang",
		Email:    "gilang@example.com",
		IsActive: true,
	}

	mockRepo.On("FindByID", mock.Anything, "user-123").Return(mockUser, nil)

	res, err := service.GetMyProfile(context.Background(), "user-123")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Gilang", res.Name)
	mockRepo.AssertExpectations(t)
}

func TestGetMyProfile_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	mockRepo.On("FindByID", mock.Anything, "invalid-id").Return(nil, nil)

	res, err := service.GetMyProfile(context.Background(), "invalid-id")

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "user not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUpdateMyProfile_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	mockUser := &User{ID: "user-123", Name: "Lama"}

	newName := "Baru"
	req := UpdateProfileRequest{Name: &newName}

	mockRepo.On("FindByID", mock.Anything, "user-123").Return(mockUser, nil)
	mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(u *User) bool {
		return u.Name == "Baru"
	})).Return(nil)

	err := service.UpdateMyProfile(context.Background(), "user-123", req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteProfilePicture_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	mockUser := &User{ID: "user-123", ProfilePicture: "foto_lama.jpg"}

	mockRepo.On("FindByID", mock.Anything, "user-123").Return(mockUser, nil)
	mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(u *User) bool {
		return u.ProfilePicture == "" // Harus kosong string
	})).Return(nil)

	err := service.DeleteProfilePicture(context.Background(), "user-123")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateMyPassword_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	// Hash password "passwordlama"
	hashed, _ := hash.MakeHash("passwordlama")
	mockUser := &User{ID: "user-123", PasswordHash: hashed}

	req := UpdatePasswordRequest{
		CurrentPassword: "passwordlama",
		NewPassword:     "passwordbaru",
		ConfirmPassword: "passwordbaru",
	}

	mockRepo.On("FindByID", mock.Anything, "user-123").Return(mockUser, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)

	err := service.UpdateMyPassword(context.Background(), "user-123", req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestVerifyEmailOTP_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	mockUser := &User{ID: "user-123", Email: "lama@example.com"}

	req := VerifyEmailOTPRequest{
		OTP:      "123456",
		NewEmail: "baru@example.com",
	}

	mockRepo.On("FindByEmail", mock.Anything, "baru@example.com").Return(nil, nil)
	mockRepo.On("FindByID", mock.Anything, "user-123").Return(mockUser, nil)
	mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(u *User) bool {
		return u.Email == "baru@example.com"
	})).Return(nil)

	err := service.VerifyEmailOTP(context.Background(), "user-123", req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestVerifyEmailOTP_EmailAlreadyUsed(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	req := VerifyEmailOTPRequest{
		OTP:      "123456",
		NewEmail: "baru@example.com",
	}

	// Simulasi email sudah dipakai orang lain
	mockRepo.On("FindByEmail", mock.Anything, "baru@example.com").Return(&User{ID: "orang-lain"}, nil)

	err := service.VerifyEmailOTP(context.Background(), "user-123", req)

	assert.Error(t, err)
	assert.Equal(t, "email is already used by another account", err.Error())
	mockRepo.AssertExpectations(t)
}

// --- PENGUJIAN MANAJEMEN KARYAWAN ---

func TestCreateManager_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	req := CreateManagerRequest{
		BranchID: 10,
		Name:     "Budi Manager",
		Email:    "budi@example.com",
		Password: "password123",
	}

	mockRepo.On("CheckBranchExists", mock.Anything, 10).Return(true, nil)
	mockRepo.On("CheckBranchManagerExists", mock.Anything, 10).Return(false, nil)
	mockRepo.On("FindByEmail", mock.Anything, "budi@example.com").Return(nil, nil)
	mockRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *User) bool {
		return u.RoleID == 2 && *u.BranchID == 10 && u.Email == "budi@example.com"
	})).Return(nil)

	err := service.CreateManager(context.Background(), req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateManager_BranchAlreadyHasManager(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	req := CreateManagerRequest{
		BranchID: 10,
		Email:    "budi@example.com",
	}

	mockRepo.On("CheckBranchExists", mock.Anything, 10).Return(true, nil)
	// Simulasi cabang sudah punya manager
	mockRepo.On("CheckBranchManagerExists", mock.Anything, 10).Return(true, nil)

	err := service.CreateManager(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "this branch already has a manager", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestCreateCashier_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	req := CreateCashierRequest{
		Name:     "Siti Kasir",
		Email:    "siti@example.com",
		Password: "password123",
	}

	managerBranchID := 5

	mockRepo.On("FindByEmail", mock.Anything, "siti@example.com").Return(nil, nil)
	mockRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *User) bool {
		return u.RoleID == 3 && *u.BranchID == 5 && u.Name == "Siti Kasir"
	})).Return(nil)

	err := service.CreateCashier(context.Background(), managerBranchID, req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetEmployees_Admin(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	mockUsers := []User{
		{ID: "m1", RoleID: 2, Name: "Manager 1"},
		{ID: "c1", RoleID: 3, Name: "Cashier 1"},
	}

	mockRepo.On("FindAllEmployees", mock.Anything).Return(mockUsers, nil)

	res, err := service.GetEmployees(context.Background(), "Admin", nil)

	assert.NoError(t, err)
	assert.Len(t, res, 2)
	mockRepo.AssertExpectations(t)
}

func TestGetEmployees_Manager(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	mockUsers := []User{
		{ID: "c1", RoleID: 3, Name: "Cashier 1"},
	}

	branchID := 5
	mockRepo.On("FindEmployeesByBranch", mock.Anything, branchID).Return(mockUsers, nil)

	res, err := service.GetEmployees(context.Background(), "Manager", &branchID)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	mockRepo.AssertExpectations(t)
}

func TestDisableEmployee_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	mockUser := &User{ID: "emp-1", IsActive: true}

	mockRepo.On("FindByID", mock.Anything, "emp-1").Return(mockUser, nil)
	mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(u *User) bool {
		return u.IsActive == false // Memastikan statusnya ditoggle
	})).Return(nil)

	err := service.DisableEmployee(context.Background(), "emp-1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
