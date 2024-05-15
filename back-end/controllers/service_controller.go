package controllers

import (
	"io/ioutil"
	"API/models"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetService
//	@Summary		Get service details
//	@Description	Get details of a specific service based on the service's name.
//	@Tags			Services
//	@Accept			json
//	@Produce		json
//	@Param			name	path	string	true	"name"
//	@Success		200	{object}	models.Service	"Service details"
//	@Failure		400	"Bad request or validation error"
//	@Failure		500	"Internal server error"
//	@Router			/services/{name} [get]
func GetService() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")

		fileData, err := ioutil.ReadFile("about.json")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var about models.About
		if err := json.Unmarshal(fileData, &about); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for _, service := range about.Server.Services {
			if service.Name == name {
				c.JSON(http.StatusOK, service)
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
	}
}
