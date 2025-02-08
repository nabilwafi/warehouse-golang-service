package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/nabilwafi/warehouse-management-system/src/config"
	"github.com/nabilwafi/warehouse-management-system/src/handlers"
	"github.com/nabilwafi/warehouse-management-system/src/repositories"
	"github.com/nabilwafi/warehouse-management-system/src/routes"
	"github.com/nabilwafi/warehouse-management-system/src/services"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}
}

func main() {
	env, err := config.NewEnv()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	db, err := config.NewDB(env.DB)
	if err != nil {
		log.Fatal("Error connection to DB")
	}

	validate := validator.New()

	r := gin.Default()

	transactionRepo := repositories.NewTransactionRepository(db.Conn)

	userRepo := repositories.NewUserRepository(db.Conn)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandlerImpl(userService)

	productRepo := repositories.NewProductRepository(db.Conn)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService, validate)

	locationRepo := repositories.NewLocationRepository(db.Conn)
	locationService := services.NewLocationService(locationRepo)
	locationHandler := handlers.NewLocationHandler(locationService, validate)

	orderRepo := repositories.NewOrderRepository(db.Conn)
	orderService := services.NewOrderService(orderRepo, productRepo, transactionRepo)
	orderHandler := handlers.NewOrderHandler(orderService, validate)

	router := routes.NewRouter(r, userHandler, productHandler, locationHandler, orderHandler)
	router.Start(env.Http.Port)
}
