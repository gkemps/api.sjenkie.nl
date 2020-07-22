package bunq

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type SessionServerRequest struct {
	Secret string `json:"secret"`
}

func (service *Service) CreateSession() (string, float64, error) {
	requestBody := SessionServerRequest{
		Secret: service.ApiKey,
	}

	bodyRaw, err := json.Marshal(requestBody)
	if err != nil {
		return "", 0, err
	}

	r, err := http.NewRequest(
		http.MethodPost,
		service.BaseUrl+"session-server",
		bytes.NewBuffer(bodyRaw),
	)
	if err != nil {
		return "", 0, err
	}

	resp, err := service.DoRequest(r)
	if err != nil {
		return "", 0, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}

	if resp.StatusCode > 201 {
		return "", 0, errors.New(string(respBody))
	}

	var response map[string]interface{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return "", 0, err
	}

	log.Printf("session server response: %+v", response)

	responses := response["Response"].([]interface{})

	token := responses[1].(map[string]interface{})["Token"].(map[string]interface{})
	userId := responses[2].(map[string]interface{})["UserPerson"].(map[string]interface{})

	return token["token"].(string), userId["id"].(float64), nil
}
