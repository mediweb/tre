package tre

import (
	"errors"
	"log"
	"os"

	"github.com/VojtechVitek/go-trello"
)

func NewClientWithOsExitOnErr() *trello.Client {
	client, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}
	return client
}

// Create a new client for Trello
func NewClient() (*trello.Client, error) {
	key, token := "", ""

	if key = os.Getenv("TRELLO_KEY"); key == "" {
		return nil, errors.New("TRELLO_KEY is not set")
	}

	if token = os.Getenv("TRELLO_TOKEN"); token == "" {
		return nil, errors.New("TRELLO_TOKEN is not set")
	}

	return trello.NewAuthClient(key, &token)
}
