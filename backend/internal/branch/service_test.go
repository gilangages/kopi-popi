package branch

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository adalah mock untuk Repository cabang
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) FindAll(ctx context.Context, includeInactive bool) ([]Branch, error) {
	args := m.Called(ctx, includeInactive)
	if args.Get(0) != nil {
		return args.Get(0).([]Branch), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) FindByID(ctx context.Context, id int) (*Branch, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*Branch), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, branch *Branch) error {
	args := m.Called(ctx, branch)
	return args.Error(0)
}

func (m *MockRepository) Update(ctx context.Context, branch *Branch) error {
	args := m.Called(ctx, branch)
	return args.Error(0)
}

func TestGetAllBranches_Admin(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)
	
	branches := []Branch{{ID: 1, Name: "A", IsActive: true}, {ID: 2, Name: "B", IsActive: false}}
	// Admin meminta semua (includeInactive = true)
	mockRepo.On("FindAll", mock.Anything, true).Return(branches, nil)
	
	res, err := service.GetAllBranches(context.Background(), "Admin", true)
	
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	mockRepo.AssertExpectations(t)
}

func TestGetAllBranches_Public(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)
	
	branches := []Branch{{ID: 1, Name: "A", IsActive: true}}
	// Public (role = "") meminta includeInactive = true, namun dipaksa jadi false di service
	mockRepo.On("FindAll", mock.Anything, false).Return(branches, nil)
	
	res, err := service.GetAllBranches(context.Background(), "", true) // Paksa true, harusnya ditimpa jadi false
	
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	mockRepo.AssertExpectations(t)
}

func TestCreateBranch_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)
	
	req := CreateBranchRequest{Name: "Cabang Test", Address: "Jalan Test"}
	
	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(b *Branch) bool {
		return b.Name == "Cabang Test" && b.Address == "Jalan Test" && b.IsActive == true
	})).Return(nil)
	
	err := service.CreateBranch(context.Background(), req)
	
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateBranch_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)
	
	existingBranch := &Branch{ID: 1, Name: "Lama", Address: "Alamat Lama", IsActive: true}
	
	newName := "Baru"
	req := UpdateBranchRequest{Name: &newName}
	
	mockRepo.On("FindByID", mock.Anything, 1).Return(existingBranch, nil)
	mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(b *Branch) bool {
		return b.Name == "Baru" && b.Address == "Alamat Lama"
	})).Return(nil)
	
	err := service.UpdateBranch(context.Background(), 1, req)
	
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteBranch_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)
	
	existingBranch := &Branch{ID: 1, Name: "Test", IsActive: true}
	
	mockRepo.On("FindByID", mock.Anything, 1).Return(existingBranch, nil)
	mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(b *Branch) bool {
		return b.IsActive == false // Soft delete
	})).Return(nil)
	
	err := service.DeleteBranch(context.Background(), 1)
	
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteBranch_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)
	
	mockRepo.On("FindByID", mock.Anything, 99).Return((*Branch)(nil), nil)
	
	err := service.DeleteBranch(context.Background(), 99)
	
	assert.Error(t, err)
	assert.Equal(t, "branch not found", err.Error())
	mockRepo.AssertExpectations(t)
}
