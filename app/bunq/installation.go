package bunq

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type PostInstallation struct {
	PublicKey string `json:"client_public_key"`
}

func (service *Service) Installation() (string, error) {
	requestBody := PostInstallation{
		PublicKey: service.getPublicKey(),
	}

	bodyRaw, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	r, err := http.NewRequest(
		http.MethodPost,
		service.BaseUrl+"v1/installation",
		bytes.NewBuffer(bodyRaw),
	)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return "", err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode > 201 {
		return "", errors.New(string(respBody))
	}

	var response map[string]interface{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return "", err
	}

	log.Printf("installation response: %+v", response)

	responses := response["Response"].([]interface{})

	token := responses[1].(map[string]interface{})["Token"].(map[string]interface{})
	//serverPublicKey := responses[2].(map[string]interface{})

	return token["token"].(string), nil
}
