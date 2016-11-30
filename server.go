package main

import "github.com/gin-gonic/gin"

func handler(c *gin.Context) {
	content := gin.H{"Hello": "World"}
	c.JSON(200, content)
}

func getAllDatabases(c *gin.Context) {
	content := gin.H{"DB": "Mongo"}
	c.JSON(200, content)
}

func main() {
	app := gin.Default()
	app.GET("/", handler)
	app.GET("/api/databases", getAllDatabases)
	app.Run(":3001")
}
