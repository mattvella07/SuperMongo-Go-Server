package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

//PageLimit represents the max number of mongodb documents to return at one time
const PageLimit = 20

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
		session.Close()
	}

	c.JSON(200, m)
}

func getCollectionsInDB(c *gin.Context) {
	m := make(map[string]string)
	session, err := mgo.Dial("mongodb://localhost:27017/")

	if err != nil {
		m["Error"] = "Can't get collections"
	} else {
		db := session.DB(c.Param("db"))
		cols, colErr := db.CollectionNames()
		if colErr == nil {
			for idx, item := range cols {
				m[string(idx)] = item
			}
		}
		session.Close()
	}

	c.JSON(200, m)
}

func find(c *gin.Context) {
	var result []bson.M
	session, err := mgo.Dial("mongodb://localhost:27017/")

	if err != nil {
		fmt.Println("Error: " + err.Error())
	} else {
		query := make(map[string]interface{})

		queryParam := c.Param("query")
		queryParam = strings.Trim(queryParam, "{ }")
		queryStr := strings.Split(queryParam, ",")
		for _, q := range queryStr {
			q = strings.Trim(q, " ")
			items := strings.Split(q, ":")
			if len(items) == 2 {
				key := strings.Trim(items[0], "\"")
				val := items[1]

				if intVal, intErr := strconv.Atoi(val); intErr == nil {
					query[key] = intVal
				} else if floatVal, floatErr := strconv.ParseFloat(val, 64); floatErr == nil {
					query[key] = floatVal
				} else {
					query[key] = strings.Trim(val, "\"")
				}
			}
		}

		db := session.DB(c.Param("db")).C(c.Param("col"))
		findErr := db.Find(query).Limit(PageLimit).All(&result)
		if findErr != nil {
			fmt.Println("Error: " + findErr.Error())
		}

		session.Close()
	}

	c.JSON(200, result)
}

func main() {
	app := gin.Default()
	app.GET("/", handler)
	app.GET("/api/databases", getAllDatabases)
	app.GET("/api/collections/:db", getCollectionsInDB)
	app.GET("/api/find/:db/:col/:query/:projection/:options", find) //Options contains Sort, Limit, and Skip
	app.Run(":3001")
}
