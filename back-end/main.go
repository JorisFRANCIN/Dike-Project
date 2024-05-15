package main

import (
	"github.com/gin-gonic/gin"
	"API/utils"
	routes "API/routes"
	"API/middleware"
	"github.com/swaggo/gin-swagger"
    "github.com/swaggo/files"

	"API/docs"
)

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey API_Token
// @in header
// @name token
// @description Insert your JWT token
func main() {
    docs.SwaggerInfo.Title = "Swagger API Authentification"
	docs.SwaggerInfo.Description = "This is the API server documentation"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router := gin.Default()
	router.Use(middleware.SetupCORS())
	router.Use(utils.LoggerMiddleware())
	
	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":8080")
}
