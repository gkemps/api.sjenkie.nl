package smartthings

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

const openFrontDoorUrl = "open-front-door"
const closeFrontDoorUrl = "close-front-door"

func (st *Service) OpenFrontDoor(c *gin.Context) {
	postBody, _ := ioutil.ReadAll(c.Request.Body)
	log.Printf("POST body: %s", string(postBody))

	openFrontDoorUrl, found := st.settings.Urls[openFrontDoorUrl]
	if !found {
		c.JSON(422, map[string]interface{}{"error": "url not found"})
		return
	}

	req, err := http.NewRequest("POST", st.settings.BaseUrl+openFrontDoorUrl, nil)
	if err != nil {
		message := fmt.Sprintf("could not create POST request: %s", err)
		c.JSON(422, map[string]interface{}{"error": message})
		return
	}

	st.ExecuteRequest(req, c)
}

func (st *Service) CloseFrontDoor(c *gin.Context) {
	postBody, _ := ioutil.ReadAll(c.Request.Body)
	log.Printf("POST body: %s", string(postBody))

	closeFrontDoorUrl, found := st.settings.Urls[closeFrontDoorUrl]
	if !found {
		c.JSON(422, map[string]interface{}{"error": "url not found"})
		return
	}

	req, err := http.NewRequest("POST", st.settings.BaseUrl+closeFrontDoorUrl, nil)
	if err != nil {
		message := fmt.Sprintf("could not create POST request: %s", err)
		c.JSON(422, map[string]interface{}{"error": message})
		return
	}

	st.ExecuteRequest(req, c)
}

func (st *Service) ExecuteRequest(req *http.Request, c *gin.Context) {
	req.Header.Set("Authorization", "Bearer "+st.settings.Bearer)

	result, err := st.client.Do(req)
	if err != nil {
		message := fmt.Sprintf("smarttings returned error: %s", err)
		c.JSON(422, map[string]interface{}{"error": message})
		return
	}

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		message := fmt.Sprintf("problem reading body: %s", err)
		c.JSON(422, map[string]interface{}{"error": message})
		return
	}

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		message := fmt.Sprintf("problem parsing body: %s", err)
		c.JSON(422, map[string]interface{}{"error": message})
		return
	}

	c.JSON(result.StatusCode, response)
}
