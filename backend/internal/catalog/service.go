package catalog

import (
	"context"
	"errors"
	"time"
)

type Service interface {
	// Category
	GetAllCategories(ctx context.Context) ([]Category, error)
	CreateCategory(ctx context.Context, req CategoryRequest) error
	UpdateCategory(ctx context.Context, id int, req CategoryRequest) error
	DeleteCategory(ctx context.Context, id int) error

	// Material
	GetAllMaterials(ctx context.Context) ([]Material, error)
	CreateMaterial(ctx context.Context, req MaterialRequest) error
	UpdateMaterial(ctx context.Context, id int, req MaterialRequest) error
	DeleteMaterial(ctx context.Context, id int) error

	// Product
	GetAllProducts(ctx context.Context, categoryID *int, search string) ([]Product, error)
	GetProductDetail(ctx context.Context, id int, role string, includeRecipe bool) (*Product, error)
	CreateProduct(ctx context.Context, req ProductRequest) error
	UpdateProduct(ctx context.Context, id int, req ProductRequest) error
	DeleteProduct(ctx context.Context, id int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

// -- Category --
func (s *service) GetAllCategories(ctx context.Context) ([]Category, error) {
	categories, err := s.repo.FindAllCategories(ctx)
	if err != nil {
		return nil, err
	}
	if categories == nil {
		categories = []Category{}
	}
	return categories, nil
}

func (s *service) CreateCategory(ctx context.Context, req CategoryRequest) error {
	cat := &Category{Name: req.Name}
	return s.repo.CreateCategory(ctx, cat)
}

func (s *service) UpdateCategory(ctx context.Context, id int, req CategoryRequest) error {
	cat, err := s.repo.FindCategoryByID(ctx, id)
	if err != nil {
		return err
	}
	if cat == nil {
		return errors.New("category not found")
	}
	cat.Name = req.Name
	return s.repo.UpdateCategory(ctx, cat)
}

func (s *service) DeleteCategory(ctx context.Context, id int) error {
	cat, err := s.repo.FindCategoryByID(ctx, id)
	if err != nil {
		return err
	}
	if cat == nil {
		return errors.New("category not found")
	}
	return s.repo.DeleteCategory(ctx, id)
}

// -- Material --
func (s *service) GetAllMaterials(ctx context.Context) ([]Material, error) {
	mats, err := s.repo.FindAllMaterials(ctx)
	if err != nil {
		return nil, err
	}
	if mats == nil {
		mats = []Material{}
	}
	return mats, nil
}

func (s *service) CreateMaterial(ctx context.Context, req MaterialRequest) error {
	mat := &Material{
		CategoryID: req.CategoryID,
		Name:       req.Name,
		Unit:       req.Unit,
		CreatedAt:  time.Now(),
	}
	return s.repo.CreateMaterial(ctx, mat)
}

func (s *service) UpdateMaterial(ctx context.Context, id int, req MaterialRequest) error {
	mat, err := s.repo.FindMaterialByID(ctx, id)
	if err != nil {
		return err
	}
	if mat == nil {
		return errors.New("material not found")
	}
	mat.CategoryID = req.CategoryID
	mat.Name = req.Name
	mat.Unit = req.Unit
	return s.repo.UpdateMaterial(ctx, mat)
}

func (s *service) DeleteMaterial(ctx context.Context, id int) error {
	mat, err := s.repo.FindMaterialByID(ctx, id)
	if err != nil {
		return err
	}
	if mat == nil {
		return errors.New("material not found")
	}
	return s.repo.DeleteMaterial(ctx, id)
}

// -- Product --
func (s *service) GetAllProducts(ctx context.Context, categoryID *int, search string) ([]Product, error) {
	products, err := s.repo.FindAllProducts(ctx, categoryID, search)
	if err != nil {
		return nil, err
	}
	if products == nil {
		products = []Product{}
	}
	// BOM dipastikan nil untuk list
	for i := range products {
		products[i].Recipe = nil
	}
	return products, nil
}

func (s *service) GetProductDetail(ctx context.Context, id int, role string, includeRecipe bool) (*Product, error) {
	product, err := s.repo.FindProductByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	// Filter hak akses resep (Hanya Admin / Manager yang bisa includeRecipe)
	if !includeRecipe || (role != "Admin" && role != "Manager") {
		product.Recipe = nil
	}

	return product, nil
}

func (s *service) CreateProduct(ctx context.Context, req ProductRequest) error {
	product := &Product{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		ImageURL:    req.ImageURL,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	var boms []ProductBOM
	for _, rb := range req.Recipe {
		boms = append(boms, ProductBOM{
			MaterialID:     rb.MaterialID,
			QuantityNeeded: rb.QuantityNeeded,
		})
	}

	return s.repo.CreateProductWithBOM(ctx, product, boms)
}

func (s *service) UpdateProduct(ctx context.Context, id int, req ProductRequest) error {
	product, err := s.repo.FindProductByID(ctx, id)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("product not found")
	}

	product.CategoryID = req.CategoryID
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.ImageURL = req.ImageURL
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}
	product.UpdatedAt = time.Now()

	var boms []ProductBOM
	for _, rb := range req.Recipe {
		boms = append(boms, ProductBOM{
			MaterialID:     rb.MaterialID,
			QuantityNeeded: rb.QuantityNeeded,
		})
	}

	return s.repo.UpdateProductWithBOM(ctx, product, boms)
}

func (s *service) DeleteProduct(ctx context.Context, id int) error {
	product, err := s.repo.FindProductByID(ctx, id)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("product not found")
	}

	return s.repo.SoftDeleteProduct(ctx, id)
}
