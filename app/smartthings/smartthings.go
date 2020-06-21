package smartthings

import "net/http"

type Settings struct {
	BaseUrl string `json:"base-url"`
	Urls    map[string]string
	Bearer  string
}

type Service struct {
	client   *http.Client
	settings *Settings
}

func NewSmartThingsService(settings *Settings) *Service {
	client := &http.Client{}

	return &Service{
		client:   client,
		settings: settings,
	}
}
