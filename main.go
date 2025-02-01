package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nabilwafi/warehouse-management-system/database"
	"github.com/nabilwafi/warehouse-management-system/exception"
	"github.com/nabilwafi/warehouse-management-system/handlers"
	"github.com/nabilwafi/warehouse-management-system/repositories"
	"github.com/nabilwafi/warehouse-management-system/routes"
	"github.com/nabilwafi/warehouse-management-system/services"
	"github.com/spf13/viper"
)

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	viper.AutomaticEnv()
}

func main() {
	db := database.Connection()
	defer db.Close()

	validator := validator.New()

	r := gin.Default()
	r.Use(exception.GlobalErrorHandler())

	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo, validator, db)
	userHandler := handlers.NewUserHandlerImpl(userService)

	productRepo := repositories.NewProductRepository()
	productService := services.NewProductService(productRepo, validator, db)
	productHandler := handlers.NewProductHandler(productService)

	locationRepo := repositories.NewLocationRepository()
	locationService := services.NewLocationService(locationRepo, validator, db)
	locationHandler := handlers.NewLocationHandler(locationService)

	orderRepo := repositories.NewOrderRepository()
	orderService := services.NewOrderService(orderRepo, productRepo, db)
	orderHandler := handlers.NewOrderHandler(orderService)

	app := routes.NewRouter(r, userHandler, productHandler, locationHandler, orderHandler)
	PORT := viper.GetString("PORT")
	app.Start(PORT)
}
