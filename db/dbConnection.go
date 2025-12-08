package db

import (
	_ "github.com/lib/pq"
	sql "database/sql"
	"os"
)

func HandleConnection() *sql.DB {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSslMode := os.Getenv("DB_SSLMODE")

	connStr := "user=" + dbUser + " dbname=" + dbName + " password=" + dbPassword + " sslmode=" + dbSslMode

	db, dbErr := sql.Open("postgres", connStr)
	if dbErr != nil {
		panic("Error when connecting to PostgreSQL")
	}
	return db
}