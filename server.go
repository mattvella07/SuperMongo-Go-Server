package main

import (
	"github.com/gin-gonic/gin"
	"labix.org/v2/mgo"
)

func handler(c *gin.Context) {
	content := gin.H{"Test": "Hi"}
	c.JSON(200, content)
}

func getAllDatabases(c *gin.Context) {
	m := make(map[string]string)
	session, err := mgo.Dial("mongodb://localhost:27017/")

	if err != nil {
		m["Error"] = "Can't connect to database"
	} else {
		dbs, dbErr := session.DatabaseNames()
		if dbErr == nil {
			for idx, item := range dbs {
				m[string(idx)] = item
			}
		}
	}

	c.JSON(200, m)
}

func main() {
	app := gin.Default()
	app.GET("/", handler)
	app.GET("/api/databases", getAllDatabases)
	app.Run(":3001")
}
