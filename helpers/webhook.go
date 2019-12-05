package helpers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Webhook struct {
	Step             string
	Success          bool
	ErrorExplanation string
}

func SendWebhook(data Webhook) {
	requestBody, err := json.Marshal(data)
	if err != nil {
		log.Println(err.Error())
		return
	}

	_, err = http.Post(os.Getenv("WEBHOOK_URL"), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println(err.Error())
		return
	}
}
