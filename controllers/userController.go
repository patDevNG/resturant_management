package controllers

import (
	"github.com/gin-gonic/gin"
)

func GetUsers() gin.HandlerFunc{
	return func (c *gin.Context)  {
		
	}
}

func GetUser() gin.HandlerFunc{
	return func(c *gin.Context) {}
}

func SignUp() gin.HandlerFunc{
	return func(c *gin.Context) {}
}

func Auth() gin.HandlerFunc{
	return func(c *gin.Context) {}
}

func HashPassword(password string) string {}

func verifyPassword(userPassword string, providesPassword string) (bool, string){}