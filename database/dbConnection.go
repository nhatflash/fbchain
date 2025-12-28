package database

import (
	sql "database/sql"
	"os"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectToPostgreSQL() (*sql.DB, error) {
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSslMode := os.Getenv("DB_SSLMODE")

	connStr := "user=" + "postgres" + " dbname=" + dbName + " password=" + dbPassword + " sslmode=" + dbSslMode

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}


func ConnectToRedisServer() *redis.Client {
	server := os.Getenv("REDIS_SERVER")

	rdb := redis.NewClient(&redis.Options{
		Addr: server,
	})
	return rdb
}


func ConnectToMongoDB() (*mongo.Client, error) {
	uri := os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, err
	}
	return client, nil
}