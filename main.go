package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/hibiken/asynq"
	env "github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nhatflash/fbchain/api"
	"github.com/nhatflash/fbchain/controller"
	"github.com/nhatflash/fbchain/database"
	_ "github.com/nhatflash/fbchain/docs"
	"github.com/nhatflash/fbchain/graph"
	"github.com/nhatflash/fbchain/helper"
	"github.com/nhatflash/fbchain/middleware"
	"github.com/nhatflash/fbchain/repository"
	"github.com/nhatflash/fbchain/routes"
	"github.com/nhatflash/fbchain/service"
	"github.com/nhatflash/fbchain/tasks"
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
		log.Fatalf("Error loading .env file: %v\n", err)
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
		log.Fatalf("Connect to PostgreSQL failed: %v\n", err)
	}

	defer db.Close()

	// Redis 
	rdb := database.ConnectToRedisServer()

	var mongodb *mongo.Client
	mongodb, err = database.ConnectToMongoDB()
	if err != nil {
		log.Fatalf("Connect to MongoDB failed: %v\n", err)
	}

	defer mongodb.Disconnect(context.TODO())
	err = database.ValidateRestaurantItemSchema(mongodb.Database("restaurants"))
	if err != nil {
		log.Fatalf("Error when validate restaurant items schema: %v\n", err)
	}
	rItemColl := mongodb.Database("restaurants").Collection("restaurant_items")
	
	// Dependency injection
	userRepository := repository.NewUserRepository(db)
	tenantRepository := repository.NewTenantRepository(db)
	restaurantRepository := repository.NewRestaurantRepository(db)
	restaurantItemRepository := repository.NewRestaurantItemRepository(rItemColl)
	restaurantTableRepository := repository.NewRestaurantTableRepository(db)
	restaurantOrderRepository := repository.NewRestaurantOrderRepositoty(db)
	restaurantPaymentRepository := repository.NewRestaurantPaymentRepository(db)
	subPackageRepository := repository.NewSubPackageRepository(db)
	orderRepository := repository.NewOrderRepository(db)
	paymentRepository := repository.NewPaymentRepository(db)

	authService := service.NewAuthService(userRepository, tenantRepository, rdb)
	userService := service.NewUserService(userRepository)
	tenantService := service.NewTenantService(tenantRepository, userRepository)
	restaurantService := service.NewRestaurantService(restaurantRepository, subPackageRepository, restaurantItemRepository, restaurantTableRepository, restaurantOrderRepository, restaurantPaymentRepository, rdb)
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
	err = tasks.CreateAdminUserIfNotExists(db)
	if err != nil {
		log.Fatalf("Error when perform initialize admin account: %v\n", err)
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

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, api.ApiResponse{
			Status: http.StatusOK,
			Message: "Server alive",
			Data: nil,
		})
	})



	// set up asynq server for clean up order
	asqClient := tasks.RegisterAsynqClient()
	defer asqClient.Close()

	srv, mux := tasks.RegisterAsynqServer(restaurantOrderRepository)

	go func() {
		if err = srv.Run(mux); err != nil {
			log.Fatalf("Asynq server start error: %v\n", err)
		}
	}()


	// setting up scheduler for 
	var scheduler *asynq.Scheduler
	scheduler, err = tasks.RegisterAsynqScheduler()
	if err != nil {
		log.Fatalf("Error when register Asynq scheduler: %v\n", err)
	}
	go func() {
		if err = scheduler.Run(); err != nil {
			log.Fatalf("Asynq scheduler start error: %v\n", err)
		}
	}()

	
	// set up go routine for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err = r.Run(port); err != nil {
			log.Fatalf("Gin server start error: %v\n", err)
		}
	}()

	<-stop
	log.Println("Shutting down...")
	srv.Shutdown()
}