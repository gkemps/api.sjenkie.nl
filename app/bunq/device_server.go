package bunq

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type PostDeviceServer struct {
	Description  string   `json:"description"`
	Secret       string   `json:"secret"`
	PermittedIps []string `json:"permitted_ips"`
}

func (service *Service) NewDevice(description string, ips []string) error {
	requestBody := PostDeviceServer{
		Description:  description,
		Secret:       service.ApiKey,
		PermittedIps: ips,
	}

	bodyRaw, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	r, err := http.NewRequest(
		http.MethodPost,
		service.BaseUrl+"device-server",
		bytes.NewBuffer(bodyRaw),
	)
	if err != nil {
		return err
	}

	res, err := service.DoRequest(r)
	if err != nil {
		return err
	}

	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	log.Printf("response %d device server: %s", res.StatusCode, string(respBody))

	return nil
}
