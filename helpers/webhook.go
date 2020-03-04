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
	ErrorExplanation string `json:"errorExplanation,omitempty"`
	Password         string `json:"password,omitempty"`
	VCenterId        string `json:"vCenterId,omitempty"`
}

func SendWebhook(data Webhook) {
	requestBody, err := json.Marshal(data)
	if err != nil {
		log.Println(err.Error())
		return
	}

	req, err := http.NewRequest("POST", os.Getenv("WEBHOOK_URL"), bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println(err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")

	authHeader := os.Getenv("WEBHOOK_AUTH_HEADER")
	authToken := os.Getenv("WEBHOOK_AUTH_TOKEN")
	if authToken != "" && authHeader != "" {
		req.Header.Set(authHeader, authToken)
	}

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err.Error())
		return
	}
}
