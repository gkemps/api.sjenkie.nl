package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
)

func main() {
	r := gin.Default()

	jsonFile, err := os.Open("../conf/auth.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	jsonBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	var accounts gin.Accounts
	err = json.Unmarshal(jsonBytes, &accounts)
	if err != nil {
		panic(err)
	}

	// Group using gin.BasicAuth() middleware
	// gin.Accounts is a shortcut for map[string]string
	authorized := r.Group("/", gin.BasicAuth(accounts))

	authorized.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(":8084") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
