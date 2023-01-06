package routes

import(
	controller "resturant-management/controllers"
	"github.com/gin-gonic/gin"

)

func UserRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
	incomingRoutes.POST("/users", controller.SignUp())
	incomingRoutes.PATCH("/users/:user_id", controller.Auth())
}