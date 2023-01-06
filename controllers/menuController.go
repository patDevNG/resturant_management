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

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func GetMenus() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		res,err:= menuCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"unable to fetch menu"})
		}
		var allMenus []bson.M
		if err = res.All(ctx, &allMenus); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allMenus)
	}
}

func GetMenu() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var menu models.Menu
		menuId := c.Param("menu_id")
		 err := menuCollection.FindOne(ctx, bson.M{"menu_id": menuId}).Decode(&menu)
		defer cancel()
		 if err != nil {
			msg := fmt.Sprintf("Error fetching menu")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		 }
		 c.JSON(http.StatusOK, menu)
	}
}

func CreateMenu() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var menu models.Menu
		 if err:= c.BindJSON(&menu); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		 } 
		 validateErr := validate.Struct(menu)
		 if validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validateErr.Error()})
			return
		 }
		 menu.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		 menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		 menu.ID = primitive.NewObjectID()
		 menu.Menu_id = menu.ID.Hex()
		 res, insertErr := menuCollection.InsertOne(ctx, menu)
		 if insertErr != nil {
			msg:= fmt.Sprintf("Menu was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		 }
		 defer cancel()
		 c.JSON(http.StatusOK, res)
	}
}

func UpdateMenu() gin.HandlerFunc{
	return func(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var menu models.Menu
	if err := c.BindJSON(&menu); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	menuId := c.Param("menu_id")
	filter := bson.M{"menu_id":menuId}

	var updateObj primitive.D
	if menu.Start_Date != nil && menu.End_Date != nil{
		if !inTimeSpan(*menu.Start_Date, *menu.End_Date, time.Now()){
			msg := "Kindly retype the time"
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			defer cancel()
			return
		}

		updateObj = append(updateObj, bson.E{"start_date", menu.Start_Date})
		updateObj = append(updateObj, bson.E{"end_date", menu.End_Date})
		if menu.Name !=""{
			updateObj = append(updateObj, bson.E{"name", menu.Name})
		}

		if menu.Category != ""{
			updateObj = append(updateObj, bson.E{"category", menu.Category})
		}

		menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", menu.Updated_at})

		upsert := true

		opt:= options.UpdateOptions{
			Upsert: &upsert,
		}
		update := bson.D{{"$set", updateObj}}

		res, err := menuCollection.UpdateOne(
			ctx,
			filter,
			update,
			&opt,
		)
		if err != nil{
			msg:= fmt.Sprintf("Error occured while updating")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}
		defer cancel()
		c.JSON(http.StatusOK, res)
	}
	}
}

func inTimeSpan(start,end, check time.Time) bool {
	return start.After(time.Now()) && end.After(start)
}