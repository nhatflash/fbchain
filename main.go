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
	_ "github.com/nhatflash/fbchain/docs"
	"github.com/nhatflash/fbchain/helper"
	"github.com/nhatflash/fbchain/middleware"
	"github.com/nhatflash/fbchain/routes"
	swgFiles "github.com/swaggo/files"
	ginSwg "github.com/swaggo/gin-swagger"
)

// @title FB Chain Management API
// @version 1.0
// @description API Documentation for FB Chain Management API - Developed by Ducking Team
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
func main() {
	envErr := env.Load(".env")
	if envErr != nil {
		log.Fatalln("Error loading .env file")
		return;
	}
	r := gin.Default()

	r.SetTrustedProxies(nil)
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.FilterConfigurer("http://localhost:5173"))

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		phoneVErr := v.RegisterValidation("phone", helper.PhoneNumberValidator)
		if phoneVErr != nil {
			panic(phoneVErr)
		}
		identityVErr := v.RegisterValidation("identity", helper.IdentityNumberValidator)
		if identityVErr != nil {
			panic(identityVErr)
		}
		nameVErr := v.RegisterValidation("name", helper.NameValidator)
		if nameVErr != nil {
			panic(nameVErr)
		}
		pCVErr := v.RegisterValidation("postalcode", helper.PostalCodeValidator)
		if pCVErr != nil {
			panic(pCVErr)
		}
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