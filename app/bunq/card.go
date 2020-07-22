package bunq

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (service *Service) ListCards() (string, error) {
	r, err := http.NewRequest(
		http.MethodGet,
		service.BaseUrl+fmt.Sprintf("user/%.f/card", service.UserId),
		nil,
	)
	if err != nil {
		return "", err
	}

	resp, err := service.DoRequest(r)
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

	return string(respBody), nil
}
