package database

import ("github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func ConnectDB() (*sqlx.DB, error) {
	connStr := "user=sergo_user password=13410285 dbname=market host=localhost port=5433 sslmode=disable"
	var err error
	db, err = sqlx.Connect("postgres", connStr)
	return db, err
}

func GetDB() *sqlx.DB {
	return db
}