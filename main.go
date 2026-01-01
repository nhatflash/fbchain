package main

import (
	"context"
	"database/sql"
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
	"github.com/nhatflash/fbchain/controller"
	"github.com/nhatflash/fbchain/database"
	_ "github.com/nhatflash/fbchain/docs"
	"github.com/nhatflash/fbchain/graph"
	"github.com/nhatflash/fbchain/helper"
	"github.com/nhatflash/fbchain/initializer"
	"github.com/nhatflash/fbchain/middleware"
	"github.com/nhatflash/fbchain/repository"
	"github.com/nhatflash/fbchain/routes"
	"github.com/nhatflash/fbchain/service"
	swgFiles "github.com/swaggo/files"
	ginSwg "github.com/swaggo/gin-swagger"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// @title FB Chain Management API
// @version 1.0
// @description API Documentation for FB Chain Management API - Developed by Ducking Team
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	var err error

	// Load .env file
	err = env.Load(".env")
	if err != nil {
		log.Fatalln("Error loading .env file")
		return;
	}

	r := gin.Default()

	// Middleware registration
	r.SetTrustedProxies(nil)
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.FilterConfigurer("http://localhost:5173"))

	// Connect to Postgres SQL database
	var db *sql.DB
	db, err = database.ConnectToPostgreSQL()
	if err != nil {
		log.Fatalln("Connect to PostgreSQL failed", err.Error())
		return
	}

	defer db.Close()

	// Redis 
	rdb := database.ConnectToRedisServer()

	var mongodb *mongo.Client
	mongodb, err = database.ConnectToMongoDB()
	if err != nil {
		log.Fatalln("Connect to MongoDB failed", err.Error())
		return
	}

	defer mongodb.Disconnect(context.TODO())
	err = database.ValidateRestaurantItemSchema(mongodb.Database("restaurants"))
	if err != nil {
		log.Fatalln("Error when validate restaurant items schema:", err)
		return
	}
	rItemColl := mongodb.Database("restaurants").Collection("restaurant_items")
	
	// Dependency injection
	userRepository := repository.NewUserRepository(db)
	tenantRepository := repository.NewTenantRepository(db)
	restaurantRepository := repository.NewRestaurantRepository(db)
	restaurantItemRepository := repository.NewRestaurantItemRepository(rItemColl)
	restaurantTableRepository := repository.NewRestaurantTableRepository(db)
	restaurantOrderRepository := repository.NewRestaurantOrderRepositoty(db)
	subPackageRepository := repository.NewSubPackageRepository(db)
	orderRepository := repository.NewOrderRepository(db)
	paymentRepository := repository.NewPaymentRepository(db)

	authService := service.NewAuthService(userRepository, tenantRepository, rdb)
	userService := service.NewUserService(userRepository)
	tenantService := service.NewTenantService(tenantRepository, userRepository)
	restaurantService := service.NewRestaurantService(restaurantRepository, subPackageRepository, restaurantItemRepository, restaurantTableRepository, restaurantOrderRepository, rdb)
	subPackageService := service.NewSubPackageService(subPackageRepository)
	orderService := service.NewOrderService(restaurantRepository, subPackageRepository, orderRepository)
	paymentService := service.NewPaymentService(paymentRepository, orderRepository)
	vnPayService := service.NewVnPayService(orderRepository)

	authController := controller.NewAuthController(authService)
	tenantController := controller.NewTenantController(tenantService, userService)
	restaurantController := controller.NewRestaurantController(userService, restaurantService, tenantService)
	subPackageController := controller.NewSubPackageController(subPackageService)
	orderController := controller.NewOrderController(orderService, tenantService, userService)
	userController := controller.NewUserController(userService)
	paymentController := controller.NewPaymentController(paymentService, vnPayService)

	// GraphQL handler
	gqlHandler := handler.New(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &graph.Resolver{
					UserService: userService,
					TenantService: tenantService,
					RestaurantService: restaurantService,
					SubPackageService: subPackageService,
					OrderService: orderService,
				},
			},
		),
	)

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


	// Validation binding registration
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


	// Initialize admin if not exist
	err = initializer.CreateAdminUserIfNotExists(db)
	if err != nil {
		log.Fatalln("Error when perform initialize admin account.")
		return
	}

	// Define routes for REST API
	routes.MainRoutes(r, authController, tenantController, subPackageController, restaurantController, orderController, userController, paymentController)
	r.GET("/swagger/*any", ginSwg.WrapHandler(swgFiles.Handler))


	// GraphQL routes
	r.GET("/playground", func(c *gin.Context) {
		playground.Handler("GraphQL", "/graphql")(c.Writer, c.Request)
	})
	r.POST("/graphql", middleware.JwtGraphQLHandler(), func(c *gin.Context) {
		gqlHandler.ServeHTTP(c.Writer, c.Request)
	})

	r.Run(port)
}