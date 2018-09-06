package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type listResponse struct {
	Name string `json:"name"`
}

const pageLimit = 20

func find(rw http.ResponseWriter, r *http.Request) {
	// /api/find/:db/:col/:query/:projection/:options

	session, err := mgo.Dial("mongodb://localhost:27017/")
	defer session.Close()
	if err != nil {
		log.Printf("ERROR: %s", err)
		rw.WriteHeader(500)
		rw.Write([]byte(fmt.Sprintf("ERROR: %s", err)))
		return
	}

	allParams := strings.Split(strings.Replace(r.URL.String(), "/api/find/", "", -1), "/")
	var dbName, colName, queryStr, projectionStr, optionsStr string
	if len(allParams) > 0 {
		dbName = allParams[0]
	}
	if len(allParams) > 1 {
		colName = allParams[1]
	}
	if len(allParams) > 2 {
		queryStr = allParams[2]
	}
	if len(allParams) > 3 {
		projectionStr = allParams[3]
	}
	if len(allParams) > 4 {
		optionsStr = allParams[4]
	}

	fmt.Println("db: ", dbName)
	fmt.Println("col: ", colName)
	fmt.Println("query: ", queryStr)
	fmt.Println("projection: ", projectionStr)
	fmt.Println("options: ", optionsStr)

	query := make(map[string]interface{})

	var res []bson.M
	db := session.DB(dbName).C(colName)
	err = db.Find(query).Limit(pageLimit).All(&res)
	if err != nil {
		log.Printf("ERROR: %s", err)
		rw.WriteHeader(500)
		rw.Write([]byte(fmt.Sprintf("ERROR: %s", err)))
		return
	}

	fmt.Println(res)

	session.Close()

	// 		query := make(map[string]interface{})

	// 		queryParam := c.Param("query")
	// 		queryParam = strings.Trim(queryParam, "{ }")
	// 		queryStr := strings.Split(queryParam, ",")
	// 		for _, q := range queryStr {
	// 			q = strings.Trim(q, " ")
	// 			items := strings.Split(q, ":")
	// 			if len(items) == 2 {
	// 				//Key and value to be added to the query
	// 				key := strings.Trim(items[0], "\"")
	// 				val := items[1]

	// 				//Determine data type of val and convert it
	// 				if key == "_id" {
	// 					query[key] = bson.ObjectIdHex(val)
	// 				} else if intVal, intErr := strconv.Atoi(val); intErr == nil {
	// 					query[key] = intVal
	// 				} else if floatVal, floatErr := strconv.ParseFloat(val, 64); floatErr == nil {
	// 					query[key] = floatVal
	// 				} else if boolVal, boolErr := strconv.ParseBool(val); boolErr == nil {
	// 					query[key] = boolVal
	// 				} else if dateVal, dateErr := time.Parse("2006-01-02", val); dateErr == nil {
	// 					query[key] = dateVal
	// 				} else {
	// 					query[key] = strings.Trim(val, "\"")
	// 				}
	// 			}
	// 		}

	// 		db := session.DB(c.Param("db")).C(c.Param("col"))
	// 		findErr := db.Find(query).Limit(PageLimit).All(&result)
	// 		if findErr != nil {
	// 			fmt.Println("Error: " + findErr.Error())
	// 		}

	// 		session.Close()

	// 	c.JSON(200, result)

}

func getCollectionsInDB(rw http.ResponseWriter, r *http.Request) {
	var returnCols []listResponse

	session, err := mgo.Dial("mongodb://localhost:27017/")
	defer session.Close()
	if err != nil {
		log.Printf("ERROR: %s", err)
		rw.WriteHeader(500)
		rw.Write([]byte(fmt.Sprintf("ERROR: %s", err)))
		return
	}

	dbName := strings.Replace(strings.Replace(r.URL.String(), "/api/collections/", "", -1), "/", "", -1)
	fmt.Println(dbName)

	db := session.DB(dbName)
	cols, err := db.CollectionNames()
	if err != nil {
		log.Printf("ERROR: %s", err)
		rw.WriteHeader(500)
		rw.Write([]byte(fmt.Sprintf("ERROR: %s", err)))
		return
	}

	for _, item := range cols {
		returnCols = append(returnCols, listResponse{item})
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(returnCols)
}

func getDatabases(rw http.ResponseWriter, r *http.Request) {
	var returnDbs []listResponse

	session, err := mgo.Dial("mongodb://localhost:27017/")
	defer session.Close()
	if err != nil {
		log.Printf("ERROR: %s", err)
		rw.WriteHeader(500)
		rw.Write([]byte(fmt.Sprintf("ERROR: %s", err)))
		return
	}

	dbs, err := session.DatabaseNames()
	if err != nil {
		log.Printf("ERROR: %s", err)
		rw.WriteHeader(500)
		rw.Write([]byte(fmt.Sprintf("ERROR: %s", err)))
		return
	}

	for _, item := range dbs {
		returnDbs = append(returnDbs, listResponse{item})
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(returnDbs)
}

func healthCheck(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/html")
	rw.Write([]byte("Success"))
}

func createServer() {
	http.HandleFunc("/", healthCheck)
	http.HandleFunc("/api/databases", getDatabases)
	http.HandleFunc("/api/collections/", getCollectionsInDB)
	http.HandleFunc("/api/find/", find)

	log.Fatal(http.ListenAndServe(":3001", nil))
}

func main() {
	createServer()
}

// //PageLimit represents the max number of mongodb documents to return at one time
// const PageLimit = 20

// func handler(c *gin.Context) {
// 	content := gin.H{"Test": "Hi"}
// 	c.JSON(200, content)
// }

// func getAllDatabases(c *gin.Context) {
// 	m := make(map[string]string)
// 	session, err := mgo.Dial("mongodb://localhost:27017/")

// 	if err != nil {
// 		m["Error"] = "Can't connect to database"
// 	} else {
// 		dbs, dbErr := session.DatabaseNames()
// 		if dbErr == nil {
// 			for idx, item := range dbs {
// 				m[string(idx)] = item
// 			}
// 		}
// 		session.Close()
// 	}

// 	c.JSON(200, m)
// }

// func getCollectionsInDB(c *gin.Context) {
// 	m := make(map[string]string)
// 	session, err := mgo.Dial("mongodb://localhost:27017/")

// 	if err != nil {
// 		m["Error"] = "Can't get collections"
// 	} else {
// 		db := session.DB(c.Param("db"))
// 		cols, colErr := db.CollectionNames()
// 		if colErr == nil {
// 			for idx, item := range cols {
// 				m[string(idx)] = item
// 			}
// 		}
// 		session.Close()
// 	}

// 	c.JSON(200, m)
// }

// func find(c *gin.Context) {
// 	var result []bson.M
// 	session, err := mgo.Dial("mongodb://localhost:27017/")

// 	if err != nil {
// 		fmt.Println("Error: " + err.Error())
// 	} else {
// 		query := make(map[string]interface{})

// 		queryParam := c.Param("query")
// 		queryParam = strings.Trim(queryParam, "{ }")
// 		queryStr := strings.Split(queryParam, ",")
// 		for _, q := range queryStr {
// 			q = strings.Trim(q, " ")
// 			items := strings.Split(q, ":")
// 			if len(items) == 2 {
// 				//Key and value to be added to the query
// 				key := strings.Trim(items[0], "\"")
// 				val := items[1]

// 				//Determine data type of val and convert it
// 				if key == "_id" {
// 					query[key] = bson.ObjectIdHex(val)
// 				} else if intVal, intErr := strconv.Atoi(val); intErr == nil {
// 					query[key] = intVal
// 				} else if floatVal, floatErr := strconv.ParseFloat(val, 64); floatErr == nil {
// 					query[key] = floatVal
// 				} else if boolVal, boolErr := strconv.ParseBool(val); boolErr == nil {
// 					query[key] = boolVal
// 				} else if dateVal, dateErr := time.Parse("2006-01-02", val); dateErr == nil {
// 					query[key] = dateVal
// 				} else {
// 					query[key] = strings.Trim(val, "\"")
// 				}
// 			}
// 		}

// 		db := session.DB(c.Param("db")).C(c.Param("col"))
// 		findErr := db.Find(query).Limit(PageLimit).All(&result)
// 		if findErr != nil {
// 			fmt.Println("Error: " + findErr.Error())
// 		}

// 		session.Close()
// 	}

// 	c.JSON(200, result)
// }

// func main() {
// 	app := gin.Default()
// 	app.GET("/", handler)
// 	app.GET("/api/databases", getAllDatabases)
// 	app.GET("/api/collections/:db", getCollectionsInDB)
// 	app.GET("/api/find/:db/:col/:query/:projection/:options", find) //Options contains Sort, Limit, and Skip
// 	app.Run(":3001")
// }
