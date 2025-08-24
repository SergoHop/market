package repository

import (
	"github.com/jmoiron/sqlx"
    "time"
	"market/internal/models"
	"log"
    
)

type ProductRepository interface {
	CreateProduct(product *models.Product) error
    GetActiveProducts() ([]models.Product, error)
    MarkAsSold(productID int) error
}

type productRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) CreateProduct(product *models.Product) error {
    product.ExpiresAt = time.Now().Add(2 * time.Hour) 
    
    err := r.db.QueryRow(`
        INSERT INTO products (name, price, description, image_path, created_at, expires_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id`,
        product.Name,
        product.Price,
        product.Description,
        product.ImagePath,
        time.Now(),
        product.ExpiresAt,
    ).Scan(&product.ID)
    
    return err
}

func (r *productRepository) GetActiveProducts() ([]models.Product, error) {
    var products []models.Product
    
    err := r.db.Select(&products, `
        SELECT id, name, price, description, 
               image_path, created_at, expires_at
        FROM products 
        WHERE expires_at > NOW()
        ORDER BY created_at DESC`)
    
    if err != nil {
        log.Printf("DB Query Error: %v\n", err)
        return nil, err
    }
    
    // Убрали обращение к products[0] при пустом списке
    if len(products) > 0 {
        log.Printf("Loaded %d products, first expires at: %v\n", 
            len(products), 
            products[0].ExpiresAt.Format("2006-01-02 15:04:05"))
    } else {
        log.Println("No active products found")
    }
    
    return products, nil
}

func (r *productRepository) MarkAsSold(id int) error {
    _, err := r.db.Exec(`
        UPDATE products 
        SET expires_at = NOW() 
        WHERE id = $1`, id)
    return err
}

