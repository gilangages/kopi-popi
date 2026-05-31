package catalog

import "time"

// Entitas Database

type Category struct {
	ID   int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name" gorm:"unique"`
}

type Material struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	CategoryID int       `json:"category_id"`
	Name       string    `json:"name"`
	Unit       string    `json:"unit"`
	CreatedAt  time.Time `json:"created_at"`
}

type Product struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	CategoryID  int       `json:"category_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Price       float64   `json:"price"`
	ImageURL    *string   `json:"image_url"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Has-Many Relationship untuk BOM
	Recipe []ProductBOM `json:"recipe,omitempty" gorm:"foreignKey:ProductID"`
}

type ProductBOM struct {
	ID             int     `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID      int     `json:"product_id"`
	MaterialID     int     `json:"material_id"`
	QuantityNeeded float64 `json:"quantity_needed"`
}

// Request DTOs

type CategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

type MaterialRequest struct {
	CategoryID int    `json:"category_id" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Unit       string `json:"unit" binding:"required"`
}

type ProductBOMRequest struct {
	MaterialID     int     `json:"material_id" binding:"required"`
	QuantityNeeded float64 `json:"quantity_needed" binding:"required"`
}

type ProductRequest struct {
	CategoryID  int                 `json:"category_id" binding:"required"`
	Name        string              `json:"name" binding:"required"`
	Description *string             `json:"description"`
	Price       float64             `json:"price" binding:"required,gt=0"`
	ImageURL    *string             `json:"image_url"`
	IsActive    *bool               `json:"is_active"`
	Recipe      []ProductBOMRequest `json:"recipe" binding:"required"`
}
