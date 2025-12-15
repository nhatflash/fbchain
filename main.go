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
	"github.com/nhatflash/fbchain/initializer"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/nhatflash/fbchain/graph"
	"github.com/vektah/gqlparser/v2/ast"
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

	graphqlHandler := handler.New(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &graph.Resolver{

				},
			},
		),
	)

	r.SetTrustedProxies(nil)
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.FilterConfigurer("http://localhost:5173"))

	graphqlHandler.AddTransport(transport.Options{})
	graphqlHandler.AddTransport(transport.GET{})
	graphqlHandler.AddTransport(transport.POST{})

	graphqlHandler.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	graphqlHandler.Use(extension.Introspection{})
	graphqlHandler.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("phone", helper.PhoneNumberValidator)
		_= v.RegisterValidation("identity", helper.IdentityNumberValidator)
		_ = v.RegisterValidation("name", helper.NameValidator)
		_ = v.RegisterValidation("postalcode", helper.PostalCodeValidator)
		_ = v.RegisterValidation("number", helper.PositiveNumberValidator)
		_ = v.RegisterValidation("price", helper.PriceValidator)
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

	initErr := initializer.CreateAdminUserIfNotExists(db)
	if initErr != nil {
		log.Fatalln("Error when perform initialize admin account.")
		return
	}
	routes.MainRoutes(r, db)
	r.GET("/swagger/*any", ginSwg.WrapHandler(swgFiles.Handler))

	r.GET("/graphql", func(c *gin.Context) {
		playground.Handler("GraphQL", "/graphql")(c.Writer, c.Request)
	})

	r.POST("/graphql", middleware.JwtAccessHandler(), func(c *gin.Context) {
		graphqlHandler.ServeHTTP(c.Writer, c.Request)
	})

	r.Run(port)
}