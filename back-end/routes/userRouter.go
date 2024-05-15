package routes

import (
	controller "API/controllers"
	"API/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	users_routes := incomingRoutes.Group("/users")
	users_routes.Use(middleware.Authenticate())
	{
		users_routes.GET("/", controller.GetUsers())
		users_routes.GET("/:user_id", controller.GetUser())
	}
	incomingRoutes.GET("/about.json", controller.ReadJSONFile())
	incomingRoutes.GET("/services/:name", controller.GetService())
}
