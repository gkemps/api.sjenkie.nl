package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gkemps/api.sjenkie.nl/app/smartthings"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	r := gin.Default()

	//auth config
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

	//smarttings
	jsonFile, err = os.Open("../conf/smartthings.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	jsonBytes, err = ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	var smartThingsSettings smartthings.Settings
	err = json.Unmarshal(jsonBytes, &smartThingsSettings)
	if err != nil {
		panic(err)
	}

	log.Printf("sts %+v", smartThingsSettings)

	st := smartthings.NewSmartThingsService(&smartThingsSettings)

	authorized.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	authorized.POST("/open-front-door", st.OpenFrontDoor)
	authorized.POST("/close-front-door", st.CloseFrontDoor)
	r.Run(":8084") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
