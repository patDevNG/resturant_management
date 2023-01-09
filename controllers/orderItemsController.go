package controllers

import (
	"context"
	"log"
	"net/http"
	"resturant-management/database"
	"resturant-management/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)
type OrderItemPack struct {
	Table_id *string
	Order_items []models.OrderItems
}

var orderItemCollection  *mongo.Collection = database.OpenCollection(database.Client, "OrderItem")

func GetOrderItems() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		res, err := orderItemCollection.Find(context.TODO(), bson.M{});
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var allOrderItems []bson.M

		if err = res.All(ctx, &res); err != nil {
			log.Fatal()
		}
		c.JSON(http.StatusOK, allOrderItems)
	}
}

func GetOrderItem() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second);
		var orderItem models.OrderItems

		orderItemId := c.Param("order_item_id");

		err:= orderCollection.FindOne(context.TODO(), bson.M{"order_item_id": orderItemId}).Decode(&orderItem)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

	}
}

func GetOrderItemsByOrder() gin.HandlerFunc{
	return func(c *gin.Context) {}
}

func ItemByOrder(id string)(OrderItems []primitive.M, err error){
	
}

func CreateOrderItem() gin.HandlerFunc{
	return func(c *gin.Context) {}
}

func UpdateOrderItem() gin.HandlerFunc{
	return func(c *gin.Context) {}
}

