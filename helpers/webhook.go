package helpers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type Webhook struct {
	Uuid             string `json:"uuid"`
	Step             string `json:"step"`
	Success          bool   `json:"success"`
	ErrorExplanation string `json:"errorExplanation,omitempty"`
	Password         string `json:"password,omitempty"`
	VCenterId        string `json:"vCenterId,omitempty"`
	LocationId       string `json:"locationId,omitempty"`
}

func SendWebhook(data Webhook, logger *log.Logger) {
	requestBody, err := json.Marshal(data)
	if err != nil {
		logger.Println(err.Error())
		return
	}
	logger.Println(string(requestBody))

	var location location

	if data.LocationId == "" {
		location = Config.Locations["DEFAULT"]
	} else {
		location = Config.Locations[data.LocationId]
	}

	req, err := http.NewRequest("POST", location.WebhookUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		logger.Println(err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")

	authHeader := location.WebhookAuthHeader
	authToken := location.WebhookAuthToken
	if authToken != "" && authHeader != "" {
		req.Header.Set(authHeader, authToken)
	}

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		logger.Println(err.Error())
		return
	}
}
