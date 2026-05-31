package catalog

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	// Category
	FindAllCategories(ctx context.Context) ([]Category, error)
	FindCategoryByID(ctx context.Context, id int) (*Category, error)
	CreateCategory(ctx context.Context, category *Category) error
	UpdateCategory(ctx context.Context, category *Category) error
	DeleteCategory(ctx context.Context, id int) error

	// Material
	FindAllMaterials(ctx context.Context) ([]Material, error)
	FindMaterialByID(ctx context.Context, id int) (*Material, error)
	CreateMaterial(ctx context.Context, material *Material) error
	UpdateMaterial(ctx context.Context, material *Material) error
	DeleteMaterial(ctx context.Context, id int) error

	// Product
	FindAllProducts(ctx context.Context, categoryID *int, search string) ([]Product, error)
	FindProductByID(ctx context.Context, id int) (*Product, error)
	CreateProductWithBOM(ctx context.Context, product *Product, boms []ProductBOM) error
	UpdateProductWithBOM(ctx context.Context, product *Product, boms []ProductBOM) error
	SoftDeleteProduct(ctx context.Context, id int) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

// -- Category --
func (r *repository) FindAllCategories(ctx context.Context) ([]Category, error) {
	var categories []Category
	err := r.db.WithContext(ctx).Find(&categories).Error
	return categories, err
}

func (r *repository) FindCategoryByID(ctx context.Context, id int) (*Category, error) {
	var category Category
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&category).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *repository) CreateCategory(ctx context.Context, category *Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *repository) UpdateCategory(ctx context.Context, category *Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *repository) DeleteCategory(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&Category{}, id).Error
}

// -- Material --
func (r *repository) FindAllMaterials(ctx context.Context) ([]Material, error) {
	var materials []Material
	err := r.db.WithContext(ctx).Find(&materials).Error
	return materials, err
}

func (r *repository) FindMaterialByID(ctx context.Context, id int) (*Material, error) {
	var material Material
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&material).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &material, nil
}

func (r *repository) CreateMaterial(ctx context.Context, material *Material) error {
	return r.db.WithContext(ctx).Create(material).Error
}

func (r *repository) UpdateMaterial(ctx context.Context, material *Material) error {
	return r.db.WithContext(ctx).Save(material).Error
}

func (r *repository) DeleteMaterial(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&Material{}, id).Error
}

// -- Product --
func (r *repository) FindAllProducts(ctx context.Context, categoryID *int, search string) ([]Product, error) {
	var products []Product
	query := r.db.WithContext(ctx).Where("is_active = ?", true)

	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	err := query.Find(&products).Error
	return products, err
}

func (r *repository) FindProductByID(ctx context.Context, id int) (*Product, error) {
	var product Product
	// Preload BOM for details
	err := r.db.WithContext(ctx).Preload("Recipe").Where("id = ?", id).First(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (r *repository) CreateProductWithBOM(ctx context.Context, product *Product, boms []ProductBOM) error {
	// GORM Transaction
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Create Product
		if err := tx.Create(product).Error; err != nil {
			return err
		}
		// 2. Assign ProductID to BOMs and create them
		for i := range boms {
			boms[i].ProductID = product.ID
		}
		if len(boms) > 0 {
			if err := tx.Create(&boms).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *repository) UpdateProductWithBOM(ctx context.Context, product *Product, boms []ProductBOM) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Update Product
		if err := tx.Save(product).Error; err != nil {
			return err
		}
		// 2. Hapus semua BOM lama
		if err := tx.Where("product_id = ?", product.ID).Delete(&ProductBOM{}).Error; err != nil {
			return err
		}
		// 3. Masukkan BOM baru
		for i := range boms {
			boms[i].ProductID = product.ID
		}
		if len(boms) > 0 {
			if err := tx.Create(&boms).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *repository) SoftDeleteProduct(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Model(&Product{}).Where("id = ?", id).Update("is_active", false).Error
}
