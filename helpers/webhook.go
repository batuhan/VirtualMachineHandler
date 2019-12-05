package helpers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Webhook struct {
	Uuid             string `json:"uuid"`
	Step             string `json:"step"`
	Success          bool   `json:"success"`
	ErrorExplanation string `json:"error_explanation"`
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
