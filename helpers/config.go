package helpers

import (
	"log"
	"os"
	"strings"
)

type dynamicConfig struct {
	GovcURL          string
	GovcUsername     string
	GovcPassword     string
	GovcDatacenter   string
	GovcDatastore    string
	GovcResourcePool string
	Gateway          string
	Nameservers      []string
	TargetDirectory  string
}

type config struct {
	HttpPort          string
	GovcInsecure      string
	WebhookUrl        string
	WebhookAuthHeader string
	WebhookAuthToken  string
	ActiveIds         []string
	DynamicConfigs    map[string]dynamicConfig
}

var Config = config{}

func InitConfig() {
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}
	Config.HttpPort = httpPort

	Config.GovcInsecure = os.Getenv("GOVC_INSECURE")

	Config.WebhookUrl = os.Getenv("WEBHOOK_URL")

	Config.WebhookAuthHeader = os.Getenv("WEBHOOK_AUTH_HEADER")

	Config.WebhookAuthToken = os.Getenv("WEBHOOK_AUTH_TOKEN")

	activeIds := strings.Split(os.Getenv("ACTIVE_IDS"), ",")
	Config.ActiveIds = activeIds

	Config.DynamicConfigs = make(map[string]dynamicConfig)
	for _, id := range activeIds {
		dynamicConfig := dynamicConfig{}
		dynamicConfig.GovcURL = os.Getenv(id + "_GOVC_URL")
		dynamicConfig.GovcUsername = os.Getenv(id + "_GOVC_USERNAME")
		dynamicConfig.GovcPassword = os.Getenv(id + "_GOVC_PASSWORD")
		dynamicConfig.GovcDatacenter = os.Getenv(id + "_GOVC_DATACENTER")
		dynamicConfig.GovcDatastore = os.Getenv(id + "_GOVC_DATASTORE")
		dynamicConfig.GovcResourcePool = os.Getenv(id + "_GOVC_RESOURCE_POOL")
		dynamicConfig.Gateway = os.Getenv(id + "_GATEWAY")
		dynamicConfig.Nameservers = strings.Split(os.Getenv(id+"_NAMESERVERS"), ",")
		dynamicConfig.TargetDirectory = os.Getenv(id + "_TARGET_DIRECTORY")
		Config.DynamicConfigs[id] = dynamicConfig
	}

	log.Println("generated global config from env variables")
}
