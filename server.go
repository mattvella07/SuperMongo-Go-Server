package main

import (
	"github.com/gin-gonic/gin"
	"labix.org/v2/mgo"
)

func handler(c *gin.Context) {
	var content gin.H
	session, err := mgo.Dial("mongodb://localhost:27017/")

	if err != nil {
		content = gin.H{"Error": err}
	} else {
		content = gin.H{"res": session}
	}

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
