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
	ginSwg "github.com/swaggo/gin-swagger"
	swgFiles "github.com/swaggo/files"
	_ "github.com/nhatflash/fbchain/docs"
)

// @title FB Chain Management API
// @version 1.0
// @description API Documentation for FB Chain Management API - Developed by Ducking Team
// @host localhost:8080
// @BasePath /api
func main() {
	envErr := env.Load(".env")
	if envErr != nil {
		log.Fatalln("Error loading .env file")
		return;
	}
	r := gin.Default()
	r.Use(middleware.ErrorHandler())

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("phone", helper.PhoneNumberValidator)
		_ = v.RegisterValidation("identity", helper.IdentityNumberValidator)
		_ = v.RegisterValidation("name", helper.NameValidator)
		_ = v.RegisterValidation("postalcode", helper.PostalCodeValidator)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	db, dbErr := pg.HandleConnection()
	if dbErr != nil {
		log.Fatalln("Connect to PostgreSQL failed", dbErr.Error())
		return
	}

	defer db.Close()

	routes.MainRoutes(r, db)
	r.GET("/swagger/*any", ginSwg.WrapHandler(swgFiles.Handler))

	r.Run(port)
}