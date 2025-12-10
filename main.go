package main

import (
	"log"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	env "github.com/joho/godotenv"
	_ "github.com/lib/pq"
	pg "github.com/nhatflash/fbchain/database"
	"github.com/nhatflash/fbchain/middleware"
	"github.com/nhatflash/fbchain/routes"
	"github.com/nhatflash/fbchain/helper"
)

func main() {
	envErr := env.Load(".env")
	if envErr != nil {
		log.Fatalln("Error loading .env file")
		return;
	}
	router := gin.Default()
	router.Use(middleware.ErrorHandler())

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("phone", helper.PhoneNumberValidator)
		_ = v.RegisterValidation("identity", helper.IdentityNumberValidator)
		_ = v.RegisterValidation("name", helper.NameValidator)
		_ = v.RegisterValidation("postalcode", helper.PostalCodeValidator)
	}

	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		serverPort = ":8080"
	}

	db, dbErr := pg.HandleConnection()
	if dbErr != nil {
		log.Fatalln("Connect to PostgreSQL failed", dbErr.Error())
		return
	}

	defer db.Close()

	routes.MainRoutes(router, db)

	router.Run(serverPort)
}