package helpers

import "os"

type config struct {
	HttpPort          string
	GovcInsecure      string
	WebhookUrl        string
	WebhookAuthHeader string
	WebhookAuthToken  string
}

var Config = config{}

func Init() {
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}
	Config.HttpPort = httpPort

	Config.GovcInsecure = os.Getenv("GOVC_INSECURE")

	Config.WebhookUrl = os.Getenv("WEBHOOK_URL")

	Config.WebhookAuthHeader = os.Getenv("WEBHOOK_AUTH_HEADER")

	Config.WebhookAuthToken = os.Getenv("WEBHOOK_AUTH_TOKEN")
}
