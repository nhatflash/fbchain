package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	routes "github.com/nhatflash/fbchain/routes"
	_ "github.com/lib/pq"
	pg "github.com/nhatflash/fbchain/database"
	"os"
	"log"
	env "github.com/joho/godotenv"
	"github.com/nhatflash/fbchain/middleware"
)

func main() {
	fmt.Println("Starting project...")

	envErr := env.Load(".env")
	if envErr != nil {
		log.Fatalln("Error loading .env file")
		return;
	}
	router := gin.Default()
	router.Use(middleware.ErrorHandler())

	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		serverPort = ":8080"
	}

	db := pg.HandleConnection()
	if db == nil {
		log.Fatalln("Internal server error when trying to connect to DB")
		return
	}

	defer db.Close()

	routes.MainRoutes(router, db)

	router.Run(serverPort)
}