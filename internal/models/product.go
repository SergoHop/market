package models

import(
	"time"
)

type Product struct {
    ID          int        `db:"id"`
    Name        string     `db:"name"`
    Price       float64    `db:"price"`
    Description string     `db:"description"`
    ImagePath   *string    `db:"image_path"`
    CreatedAt   time.Time  `db:"created_at"`
    ExpiresAt   time.Time  `db:"expires_at"`
}
