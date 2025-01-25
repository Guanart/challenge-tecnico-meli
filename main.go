package main

import (
	"image-vuln-scanner-api/db"
	"image-vuln-scanner-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	db.ConnectDatabase()
	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.GET("images", getImages)
		v1.GET("image/:name", getImageByName)
		v1.POST("image", addImage)
	}
	router.Run(":8080")
}

func getImages(c *gin.Context) { // c *gin.Context es el contexto de la petición HTTP. Tiene información sobre la petición y la respuesta. Podemos modificar el estado de Gin desde este contexto.
	images, err := models.GetImages()
	db.CheckError(err)

	if images == nil {
		c.JSON(http.StatusNotFound, gin.H{"data": images, "message": "Resource not found"})
		return
	} else {
		var message string
		if len(images) > 0 {
			message = "Images found"
		} else {
			message = "Images not found"
		}
		c.JSON(http.StatusOK, gin.H{"data": images, "message": message})
		return
	}
}

func getImageByName(c *gin.Context) {
	name := c.Param("name")
	c.JSON(http.StatusOK, gin.H{"message": "getImageByName " + name + " Called"})
}

func addImage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "addImage Called"})
}
