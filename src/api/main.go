package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/matiasgarcia78/bgh-home-bridge/src/api/solidmation"
)

var api = solidmation.NewSolidmationApi(solidmation.Auth{User: "XXXX", Password: "XXXX"})

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/bgh/home/bridge/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"value":  api.GetStatus()})
		return
	})

	r.GET("/bgh/home/bridge/device_status", func(c *gin.Context) {
		e := c.DefaultQuery("device_id", "0")
		deviceID, err := strconv.ParseUint(e, 0, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": http.StatusText(http.StatusBadRequest), "message": fmt.Sprintf("invalid %v device id", e)})
			return
		}
		if deviceID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": http.StatusText(http.StatusBadRequest), "message": "device id is required"})
		}

		t := c.DefaultQuery("type_id", "0")
		tID, err := strconv.ParseUint(t, 0, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": http.StatusText(http.StatusBadRequest), "message": fmt.Sprintf("invalid %v type id", t)})
			return
		}

		value, err := api.GetDeviceStatus(solidmation.DeviceID(deviceID), tID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": http.StatusText(http.StatusInternalServerError), "message": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"value": value})
	})

	r.POST("/bgh/home/bridge/device_status", func(c *gin.Context) {
		t := c.DefaultQuery("temperature", "0")
		temperature, err := strconv.ParseUint(t, 0, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": http.StatusText(http.StatusBadRequest), "message": fmt.Sprintf("invalid %v temperature", t)})
			return
		}

		mode := c.DefaultQuery("mode", "")
		if len(mode) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": http.StatusText(http.StatusBadRequest), "message": "mode is required"})
			return
		}

		e := c.DefaultQuery("device_id", "0")
		deviceID, err := strconv.ParseUint(e, 0, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": http.StatusText(http.StatusBadRequest), "message": fmt.Sprintf("invalid %v device id", e)})
			return
		}
		if deviceID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": http.StatusText(http.StatusBadRequest), "message": "device id is required"})
		}

		if err := api.SetDeviceStatus(solidmation.DeviceID(deviceID), temperature, mode); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": http.StatusText(http.StatusInternalServerError), "message": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "set device mode successfully!"})

	})

	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
