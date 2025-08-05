package models

import(
	"time"
)

type Product struct {
    ID          int        `db:"id" json:"id"`
    Name        string     `db:"name" json:"name"`
    Price       float64    `db:"price" json:"price"`
    Description string     `db:"description" json:"description"`
    ImagePath   *string    `db:"image_path" json:"image_path"`
    CreatedAt   time.Time  `db:"created_at" json:"created_at"`
    ExpiresAt   time.Time  `db:"expires_at" json:"expires_at"`
}
