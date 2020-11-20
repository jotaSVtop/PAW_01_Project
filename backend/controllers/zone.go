package controllers

import (
	"net/http"
	"projetoapi/model"
	"projetoapi/services"

	"github.com/gin-gonic/gin"
)

func GetZones(c *gin.Context) {
	var zones []model.Zone

	services.Db.Find(&zones)

	if len(zones) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Didn't find any zone!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": zones})
}

func AddZone(c *gin.Context) {
	var zone model.Zone

	if err := c.ShouldBindJSON(&zone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Check syntax!"})
		return
	}
	services.Db.Save(&zone)
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Create successful!", "resourceId": zone.ID})
}

func DeleteZone(c *gin.Context) {
	var zone model.Zone

	id := c.Param("id")
	services.Db.First(&zone, id)

	if zone.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "None found!"})
		return
	}

	services.Db.Delete(&zone)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Delete succeeded!"})
}