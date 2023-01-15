package routes

import (
	"resturant-management/controllers"

	"github.com/gin-gonic/gin"
)

func EntryRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/",controllers.EntryHandler())
}