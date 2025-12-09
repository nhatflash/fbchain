package database

import (
	_ "github.com/lib/pq"
	sql "database/sql"
	"os"
)

func HandleConnection() (*sql.DB, error) {
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSslMode := os.Getenv("DB_SSLMODE")

	connStr := "user=" + "postgres" + " dbname=" + dbName + " password=" + dbPassword + " sslmode=" + dbSslMode

	db, dbErr := sql.Open("postgres", connStr)
	if dbErr != nil {
		return nil, dbErr
	}
	return db, nil
}