package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	routes "github.com/nhatflash/fbchain/routes"
	_ "github.com/lib/pq"
	pg "github.com/nhatflash/fbchain/db"
	"os"
	"log"
	env "github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting project...")

	envErr := env.Load(".env")
	if envErr != nil {
		log.Fatalln("Error loading .env file")
		return;
	}
	router := gin.Default()

	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		serverPort = ":8080"
	}

	db := pg.HandleConnection()

	defer db.Close()

	routes.MainRoutes(router)

	router.Run(serverPort)
}