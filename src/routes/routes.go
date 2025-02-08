package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nabilwafi/warehouse-management-system/src/handlers"
	"github.com/nabilwafi/warehouse-management-system/src/middlewares"
)

type router struct {
	router *gin.Engine

	user     handlers.UserHandler
	product  handlers.ProductHandler
	location handlers.LocationHandler
	order    handlers.OrderHandler
}

func NewRouter(r *gin.Engine, user handlers.UserHandler, product handlers.ProductHandler, location handlers.LocationHandler, order handlers.OrderHandler) *router {
	return &router{
		router:   r,
		user:     user,
		product:  product,
		location: location,
		order:    order,
	}
}

func (r *router) Start(port string) {
	v1 := r.router.Group("/api/v1")
	{
		v1.Use(middlewares.JWTMiddleware())

		v1.POST("/register", r.user.Register)
		v1.POST("/login", r.user.Login)

		users := v1.Group("/users")
		{
			users.GET("/me", middlewares.RoleMiddleware("staff", "admin"), r.user.GetMe)
			users.GET("/", middlewares.RoleMiddleware("admin"), r.user.ListUsers)
		}

		products := v1.Group("/products")
		{
			products.POST("/", middlewares.RoleMiddleware("admin"), r.product.AddProduct)
			products.GET("/", middlewares.RoleMiddleware("staff", "admin"), r.product.GetAllProducts)
			products.GET("/:product_id", middlewares.RoleMiddleware("staff", "admin"), r.product.GetProductByID)
			products.PUT("/:product_id", middlewares.RoleMiddleware("admin"), r.product.UpdateProduct)
			products.DELETE("/:product_id", middlewares.RoleMiddleware("admin"), r.product.DeleteProduct)
		}

		location := v1.Group("/locations")
		{
			location.POST("/", middlewares.RoleMiddleware("admin"), r.location.AddLocation)
			location.GET("/", middlewares.RoleMiddleware("staff", "admin"), r.location.GetAllLocations)
		}

		orders := v1.Group("/orders")
		{
			orders.POST("/receive", middlewares.RoleMiddleware("staff"), r.order.ReceiveOrder)
			orders.POST("/ship", middlewares.RoleMiddleware("staff"), r.order.ShipOrder)
			orders.GET("/", middlewares.RoleMiddleware("admin", "staff"), r.order.GetAllOrders)
			orders.GET("/:order_id", middlewares.RoleMiddleware("admin", "staff"), r.order.GetOrderByID)
		}
	}

	r.router.Run(fmt.Sprintf(":%s", port))
}
