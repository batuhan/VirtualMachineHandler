package helpers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type Webhook struct {
	Step             string
	Success          bool
	ErrorExplanation string
}

func SendWebhook(data Webhook) error {
	requestBody, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = http.Post(os.Getenv("WEBHOOK_URL"), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	return nil
}
