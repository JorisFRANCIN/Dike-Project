package controllers

import (
	"time"
	"io/ioutil"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)


// about.json
//	@Summary		Get services details
//	@Description	Get details of the services available as well as their actions and reactions.
//	@Accept			json
//	@Produce		json
//	@Success		200	"Services details, host IP address and current time in the Epoch Unix Time Stamp format"
//	@Failure		400	"Bad request or validation error"
//	@Failure		500	"Internal server error"
//	@Router			/about.json [get]
func ReadJSONFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		aboutJSON, err := readAboutJSONFile()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		clientHost := c.ClientIP()
		currentTime := getCurrentTimeInSeconds()
		aboutJSON["client"].(map[string]interface{})["host"] = clientHost
		aboutJSON["server"].(map[string]interface{})["current_time"] = currentTime
		if err := writeAboutJSONFile(aboutJSON); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, aboutJSON)
	}
}

func getCurrentTimeInSeconds() int64 {
	return time.Now().Unix()
}

func readAboutJSONFile() (map[string]interface{}, error) {
	fileData, err := ioutil.ReadFile("about.json")
	if err != nil {
		return nil, err
	}

	var aboutJSON map[string]interface{}
	if err := json.Unmarshal(fileData, &aboutJSON); err != nil {
		return nil, err
	}

	return aboutJSON, nil
}

func writeAboutJSONFile(data map[string]interface{}) error {
	fileData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("about.json", fileData, 0644)
	if err != nil {
		return err
	}

	return nil
}
