package main

import (
	"os"
	"github.com/gin-gonic/gin"
	"resturant-management/database"
	routes "resturant-management/routes"
	middlewares "resturant-management/middlewares"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection  *mongo.Collection = database.OpenCollection(database.Client, "food")
func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router:=gin.New()
	router.Use((gin.Logger()))
	routes.UserRoutes(router)
	router.Use(middlewares.Authentication())

	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.InvoiceRoutes(router)
	routes.TableRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemsRoutes(router)
	routes.EntryRoutes(router)

	router.Run(":" +port)
}
