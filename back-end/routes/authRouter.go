package routes

import (
	controller "API/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes (incomingRoutes *gin.Engine) {
	auth_routes := incomingRoutes.Group("/users")
	{
		auth_routes.POST("/register", controller.Register())
		auth_routes.POST("/login", controller.Login())
	}
}
