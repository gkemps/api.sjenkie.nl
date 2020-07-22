package bunq_test

import (
	"github.com/gkemps/api.sjenkie.nl/app/bunq"
	"log"
	"testing"
)

func TestListCard(t *testing.T) {
	bunqService := getBunqService()

	res, err := bunqService.ListCards()
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("list card response: %+v", res)
}

func getBunqService() *bunq.Service {
	service, err := bunq.NewService(
		"https://public-api.sandbox.bunq.com/v1/",
		"sandbox_3b24855f212fb7fb6e1e50f6a8a5d6392c0e2e6f97cd7e8d93fc3a52",
	)
	if err != nil {
		panic(err)
	}

	return service
}
