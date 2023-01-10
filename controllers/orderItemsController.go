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
	"go.mongodb.org/mongo-driver/mongo/options"
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

		err:= orderCollection.FindOne(ctx, bson.M{"order_item_id": orderItemId}).Decode(&orderItem)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, orderItem)

	}
}

func GetOrderItemsByOrder() gin.HandlerFunc{
	return func(c *gin.Context) {
		orderId := c.Param("order_id")

		allOrderItems, err :=ItemsByOrder(orderId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Error while listing order items bu order id"})
		}

		c.JSON(http.StatusOK, allOrderItems)
	}
}

func ItemsByOrder(id string)(OrderItems []primitive.M, err error){
	
}

func CreateOrderItem() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var orderItemPack OrderItemPack
		var order models.Order

		if err:= c.BindJSON(&orderItemPack); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		order.Order_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		orderItemsToBeInserted := []interface{}{}
		order_id := OrderItemOrderCreator(order)

		for _, orderItem :=range orderItemPack.Order_items{
			orderItem.Order_id = order_id

			validationErr := validate.Struct(orderItem)
			if validationErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
				return
			}
			orderItem.ID = primitive.NewObjectID()
			orderItem.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339)) 
			orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Order_items_id = orderItem.ID.Hex()
			var num = toFixed(*orderItem.Unit_price, 2)
			orderItem.Unit_price = &num
			orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)
		}
		insertedOrderItems, err := orderItemCollection.InsertMany(ctx,orderItemsToBeInserted)

		if err != nil{
			log.Fatal(err)
		}
		defer cancel()
		c.JSON(http.StatusOK, insertedOrderItems)
	}
}

func UpdateOrderItem() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var orderItem models.OrderItems
    	orderItemId := c.Param("order_item_id")
		filter := bson.M{"order_item_id": orderItemId}
		var updateObj primitive.D
		if orderItem.Unit_price !=nil {
			updateObj = append(updateObj, bson.E{"unit_price", *&orderItem.Unit_price})
		}
		
		if orderItem.Quantity != nil {
			updateObj = append(updateObj, bson.E{"quantity", orderItem.Quantity})
		}

		if orderItem.Food_id != nil {
			updateObj = append(updateObj, bson.E{"food_id",orderItem.Food_id})
		}

		if orderItem.Unit_price != nil {
			updateObj = append(updateObj, bson.E{"unit_price",orderItem.Unit_price})
		}
		orderItem.Updated_at,_ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj =append(updateObj, bson.E{"updated_at", orderItem.Updated_at})

		upsert := true
		opt:= options.UpdateOptions{
			Upsert: &upsert,
		}

		res, err := orderItemCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)
		if err != nil {
			msg:= "order item fialed"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}
		defer cancel()
		c.JSON(http.StatusOK, res)
	}
}

