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
	router.Run(":8081")
}

func getImages(c *gin.Context) { // c *gin.Context es el contexto de la petición HTTP. Tiene información sobre la petición y la respuesta. Podemos modificar el estado de Gin desde este contexto.
	images, err := models.GetImages()
	db.CheckError(err)

	if images == nil {
		c.JSON(http.StatusNotFound, gin.H{"data": images, "error": "Resource not found"})
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
	image, err := models.GetImageByName(name)
	db.CheckError(err)

	if image.Name == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": image, "message": "Image found"})
		return
	}
}

func addImage(c *gin.Context) {
	var json models.Image // Declaramos una variable de tipo Image

	if err := c.ShouldBindJSON(&json); err != nil { // ShoulBindJSON recibe el puntero a una variable, y la rellena con los datos del JSON que recibe en la petición. Si hay un error, lo devuelve
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := models.AddImage(json)

	if success {
		// Ejecutar un hilo que escanee la imagen
		go ScanImage(json.Name)
		c.JSON(http.StatusOK, gin.H{"message": "Image added successfully. Scanning image..."})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
