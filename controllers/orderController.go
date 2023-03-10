package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"resturant-management/database"
	"resturant-management/models"
	"time"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
var orderCollection *mongo.Collection = database.OpenCollection(database.Client,"order")

func GetOrders() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		 res,err := orderCollection.Find(context.TODO(), bson.M{})
		 defer cancel()

		 if err !=nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Error fetching orders"})
		 }
		
		 var allOrders []bson.M
		 
		 if err = res.All(ctx, &allOrders); err != nil{
			log.Fatal(err)
		 }

		 c.JSON(http.StatusOK, allOrders)
	}
}

func GetOrder() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second);
		orderId := c.Param("order_id")
		var order models.Order;
		err := orderCollection.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&order)

		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Error fetching order"})
		}
	}	
}

func CreateOrder() gin.HandlerFunc{
	return func(c *gin.Context) {
	var	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var table models.Table
	var order models.Order

	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validationErr := validate.Struct(order)

	if validationErr != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
	}
	
	if order.Table_id != nil {
		err := tableCollection.FindOne(ctx, bson.M{"table_id": order.Table_id}).Decode((&table))
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.Order_id = order.ID.Hex()

	res, insertErr := orderCollection.InsertOne(ctx, order)

	if insertErr != nil {
		msg:= fmt.Sprintf("order not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
	}

	defer cancel()
	c.JSON(http.StatusOK, res)
	}
}

func UpdateOrder() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var table models.Table
		var order models.Order

		orderId := c.Param("order_id")

		var updateObj primitive.D

		if err:= c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}
		if order.Table_id != nil {
			err := tableCollection.FindOne(ctx, bson.M{"table_id":order.Table_id}).Decode(&table)
			defer cancel()
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
				return
			}
			updateObj = append(updateObj, bson.E{"table_id", order.Table_id})
		}
		order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		updateObj = append(updateObj, bson.E{"updated_at", order.Updated_at})

		upsert := true

		filter := bson.M{"order_id":orderId}

		opt:= options.UpdateOptions{
			Upsert: &upsert,
		}

		update := bson.D{{"$set", updateObj}}

		res, err := orderCollection.UpdateOne(
			ctx,
			filter,
			update,
			&opt,
		)
		defer cancel()
		if err != nil {
			msg := fmt.Sprintf("Error occured while updating order")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return 
		}
		c.JSON(http.StatusOK, res)
	}
}

func OrderItemOrderCreator(order models.Order) string {
	var	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	
	order.ID = primitive.NewObjectID()
	order.Order_id = order.ID.Hex()
	
	orderCollection.InsertOne(ctx, order)
	defer cancel()
	return order.Order_id
}