package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"resturant-management/database"
	"resturant-management/helpers"
	"resturant-management/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)
var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
func GetUsers() gin.HandlerFunc{
	return func (c *gin.Context)  {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		recordPerPage, err := strconv.Atoi((c.Query("recordPerPage")))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, pageErr := strconv.Atoi(c.Query("page"))
		if pageErr != nil || page < 1 {
			page = 1
		}
		startIndex := (page-1)* recordPerPage
        startIndex, indexErr = strconv.Atoi(c.Query("startIndex"))

		matchStage:= bson.D{{"$match", bson.D{{}}}}

		projectStage:= bson.D{
			{
				"$project", bson.D{
					{"_id",0},
					{"total_count", 1},
					{"user_items", bson.D{{"$slice",[]interface{}{"data", startIndex,recordPerPage}}}},
				},
			},
		}

		res, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, projectStage})
			defer cancel()
			if err != nil{
				c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while listing items"})
				return
			}
			var allUsers []bson.M

			if err = res.All(ctx, &allUsers); err != nil{
				log.Fatal(err)
			}
			c.JSON(http.StatusOK, allUsers[0])
	}
}

func GetUser() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		userId := c.Param("user_id")
		var user models.User

		err:= userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, user)
	}
}

func SignUp() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, canel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		if err := c.BindJSON(&user); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := validate.Struct(user)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email":user.Email})
		defer canel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while checking email"})
			return
		}
		password :=HashPassword(*&user.Password)
		user.Password = *&password
		count, phoneErr := userCollection.CountDocuments(ctx,bson.M{"phone":user.Phone} )
		if phoneErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while checking phone"})
			return
		}
		defer canel()

		if count> 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"email or phone already exist"})
			return
		}

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		token, refreshToken, _ = helpers.GenerateAllTokens(*&user.Email, *&user.First_name, *&user.Last_name, user.User_id)

		user.Token = &token
		user.Refresh_Token = &refreshToken

		insertRes, insertErr := userCollection.InsertOne(ctx,user)
		if insertErr != nil {
			msg:=fmt.Sprint("user was no created")
			c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
		}
		defer canel()
		c.JSON(http.StatusOK, insertRes)
	}
}

func Auth() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err:= userCollection.FindOne(ctx, bson.M{"email":user.Email}).Decode(&foundUser)
		defer cancel()

		if err!= nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"user not found"})
		}
		passwordIsValid, msg := verifyPassword(*&user.Password, foundUser.Password)
		defer cancel()
		if passwordIsValid != true{
			c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
			return
		}
		token, refreshToken, _ = helpers.GenerateAllTokens(*&foundUser.Email, *&foundUser.First_name, *&foundUser.Last_name, *&foundUser.User_id)

		helpers.UpdateAllTokens(token, refreshToken, foundUser.User_id)
		c.JSON(http.StatusOK, foundUser)
	}
}

func HashPassword(password string) string {
	bytes,err:= bcrypt.GenerateFromPassword([]byte(password),14)
	if err!= nil{
		log.Panic(err)
	}
	return string(bytes)
}

func verifyPassword(userPassword string, providesPassword string) (bool, string){
	err:= bcrypt.CompareHashAndPassword([]byte(providesPassword),[]byte(userPassword))
	check:= true
	msg:=""
	if err !=nil{
		msg = fmt.Sprintf("login or password is incorrect")
		check = false
	}
	return check,msg
}