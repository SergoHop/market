package repository

import (
	"github.com/jmoiron/sqlx"
	_"time"
	"market/internal/models"
	_"log"
    "fmt"
)

type ProductRepository interface {
	CreateProduct(product *models.Product) error
    GetActiveProducts() ([]models.Product, error)
}

type productRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) CreateProduct(product *models.Product) error {
	query := `
		INSERT INTO products 
		(name, price, description, created_at, expires_at) 
		VALUES ($1, $2, $3, NOW(), NOW() + INTERVAL '2 hours') 
		RETURNING id`
	
	return r.db.QueryRow(
		query,
		product.Name,
		product.Price,
		product.Description,
	).Scan(&product.ID)
}

func (r *productRepository) GetActiveProducts() ([]models.Product, error) {
    var products []models.Product
    err := r.db.Select(&products, `
        SELECT 
            id, 
            name, 
            price, 
            description, 
            image_path,
            created_at,
            expires_at
        FROM products 
        WHERE expires_at > NOW()
        ORDER BY created_at DESC`)
    
    if err != nil {
        return nil, fmt.Errorf("ошибка запроса товаров: %w", err)
    }
    return products, nil
}

