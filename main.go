package main

import (
	"context"
	"log"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	env "github.com/joho/godotenv"
	_ "github.com/lib/pq"
	pg "github.com/nhatflash/fbchain/database"
	_ "github.com/nhatflash/fbchain/docs"
	"github.com/nhatflash/fbchain/graph"
	"github.com/nhatflash/fbchain/helper"
	"github.com/nhatflash/fbchain/initializer"
	"github.com/nhatflash/fbchain/middleware"
	"github.com/nhatflash/fbchain/routes"
	"github.com/nhatflash/fbchain/service"
	swgFiles "github.com/swaggo/files"
	ginSwg "github.com/swaggo/gin-swagger"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
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

	db, dbErr := pg.ConnectToDatabase()
	if dbErr != nil {
		log.Fatalln("Connect to PostgreSQL failed", dbErr.Error())
		return
	}

	defer db.Close()

	userService := service.NewUserService(db)
	tenantService := service.NewTenantService(db)

	gqlHandler := handler.New(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &graph.Resolver{
					UserService: userService,
					TenantService: tenantService,
				},
			},
		),
	)

	r.SetTrustedProxies(nil)
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.FilterConfigurer("http://localhost:5173"))

	gqlHandler.AddTransport(transport.Options{})
	gqlHandler.AddTransport(transport.GET{})
	gqlHandler.AddTransport(transport.POST{})

	gqlHandler.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	gqlHandler.Use(extension.Introspection{})
	gqlHandler.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})
	gqlHandler.SetErrorPresenter(func (ctx context.Context, err error) *gqlerror.Error {
		return gqlerror.Errorf(err.Error(), "Error code: 404")
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

	r.POST("/graphql", middleware.JwtGraphQLHandler(), func(c *gin.Context) {
		gqlHandler.ServeHTTP(c.Writer, c.Request)
	})

	r.Run(port)
}