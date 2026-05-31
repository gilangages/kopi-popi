package catalog

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

// -- Category --
func (m *MockRepository) FindAllCategories(ctx context.Context) ([]Category, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]Category), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockRepository) FindCategoryByID(ctx context.Context, id int) (*Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*Category), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockRepository) CreateCategory(ctx context.Context, category *Category) error {
	return m.Called(ctx, category).Error(0)
}
func (m *MockRepository) UpdateCategory(ctx context.Context, category *Category) error {
	return m.Called(ctx, category).Error(0)
}
func (m *MockRepository) DeleteCategory(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}

// -- Material --
func (m *MockRepository) FindAllMaterials(ctx context.Context) ([]Material, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]Material), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockRepository) FindMaterialByID(ctx context.Context, id int) (*Material, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*Material), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockRepository) CreateMaterial(ctx context.Context, material *Material) error {
	return m.Called(ctx, material).Error(0)
}
func (m *MockRepository) UpdateMaterial(ctx context.Context, material *Material) error {
	return m.Called(ctx, material).Error(0)
}
func (m *MockRepository) DeleteMaterial(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}

// -- Product --
func (m *MockRepository) FindAllProducts(ctx context.Context, categoryID *int, search string) ([]Product, error) {
	args := m.Called(ctx, categoryID, search)
	if args.Get(0) != nil {
		return args.Get(0).([]Product), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockRepository) FindProductByID(ctx context.Context, id int) (*Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*Product), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockRepository) CreateProductWithBOM(ctx context.Context, product *Product, boms []ProductBOM) error {
	return m.Called(ctx, product, boms).Error(0)
}
func (m *MockRepository) UpdateProductWithBOM(ctx context.Context, product *Product, boms []ProductBOM) error {
	return m.Called(ctx, product, boms).Error(0)
}
func (m *MockRepository) SoftDeleteProduct(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}


func TestGetProductDetail_Admin_IncludeRecipe(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	product := &Product{
		ID:   1,
		Name: "Kopi",
		Recipe: []ProductBOM{
			{MaterialID: 1, QuantityNeeded: 10},
		},
	}

	mockRepo.On("FindProductByID", mock.Anything, 1).Return(product, nil)

	res, err := service.GetProductDetail(context.Background(), 1, "Admin", true)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Recipe, 1)
	mockRepo.AssertExpectations(t)
}

func TestGetProductDetail_Cashier_TryIncludeRecipe(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	product := &Product{
		ID:   1,
		Name: "Kopi",
		Recipe: []ProductBOM{
			{MaterialID: 1, QuantityNeeded: 10},
		},
	}

	mockRepo.On("FindProductByID", mock.Anything, 1).Return(product, nil)

	// Cashier tries to include recipe, should be blocked
	res, err := service.GetProductDetail(context.Background(), 1, "Cashier", true)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, res.Recipe) // Recipe should be nil
	mockRepo.AssertExpectations(t)
}

func TestCreateProduct_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	req := ProductRequest{
		CategoryID: 1,
		Name:       "Test",
		Price:      10000,
		Recipe: []ProductBOMRequest{
			{MaterialID: 1, QuantityNeeded: 50},
		},
	}

	mockRepo.On("CreateProductWithBOM", mock.Anything, mock.MatchedBy(func(p *Product) bool {
		return p.Name == "Test" && p.Price == 10000 && p.IsActive == true
	}), mock.MatchedBy(func(boms []ProductBOM) bool {
		return len(boms) == 1 && boms[0].QuantityNeeded == 50
	})).Return(nil)

	err := service.CreateProduct(context.Background(), req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateProduct_DBError(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	req := ProductRequest{Name: "Test"}

	mockRepo.On("CreateProductWithBOM", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("db error"))

	err := service.CreateProduct(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
	mockRepo.AssertExpectations(t)
}
