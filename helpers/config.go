package helpers

import (
	"log"
	"os"
	"strings"
)

type location struct {
	GovcInsecure     string // @TODO: probably should be a boolean
	GovcURL          string
	GovcUsername     string
	GovcPassword     string
	GovcDatacenter   string
	GovcDatastore    string
	GovcResourcePool string

	Gateway         string
	Nameservers     []string
	TargetDirectory string
	DeleteDirectory  string

	WebhookUrl        string
	WebhookAuthHeader string
	WebhookAuthToken  string
}

type config struct {
	HttpPort  string
	PowerOffTimeout   string
	Locations map[string]location

}

var Config = config{}

func getLocationConfigFromEnv(locationId string, defaults location) location {
	location := location{}

	location.GovcInsecure = os.Getenv(locationId + "_GOVC_INSECURE")
	if location.GovcInsecure == "" {
		location.GovcInsecure = defaults.GovcInsecure
	}

	location.GovcURL = os.Getenv(locationId + "_GOVC_URL")
	if location.GovcURL == "" {
		location.GovcURL = defaults.GovcURL
	}

	location.GovcUsername = os.Getenv(locationId + "_GOVC_USERNAME")
	if location.GovcUsername == "" {
		location.GovcUsername = defaults.GovcUsername
	}

	location.GovcPassword = os.Getenv(locationId + "_GOVC_PASSWORD")
	if location.GovcPassword == "" {
		location.GovcPassword = defaults.GovcPassword
	}

	location.GovcDatacenter = os.Getenv(locationId + "_GOVC_DATACENTER")
	if location.GovcDatacenter == "" {
		location.GovcDatacenter = defaults.GovcDatacenter
	}

	location.GovcDatastore = os.Getenv(locationId + "_GOVC_DATASTORE")
	if location.GovcDatastore == "" {
		location.GovcDatastore = defaults.GovcDatastore
	}

	location.GovcResourcePool = os.Getenv(locationId + "_GOVC_RESOURCE_POOL")
	if location.GovcResourcePool == "" {
		location.GovcResourcePool = defaults.GovcResourcePool
	}

	location.Gateway = os.Getenv(locationId + "_GATEWAY")
	if location.Gateway == "" {
		location.Gateway = defaults.Gateway
	}

	location.Nameservers = strings.Split(os.Getenv(locationId+"_NAMESERVERS"), ",")
	if len(location.Nameservers) == 0 {
		location.Nameservers = defaults.Nameservers
	}

	location.TargetDirectory = os.Getenv(locationId + "_TARGET_DIRECTORY")
	if location.TargetDirectory == "" {
		location.TargetDirectory = defaults.TargetDirectory
	}

	location.DeleteDirectory = os.Getenv(locationId + "_DELETE_DIRECTORY")
	if location.DeleteDirectory == "" {
		location.DeleteDirectory = defaults.DeleteDirectory
	}

	location.WebhookUrl = os.Getenv("WEBHOOK_URL")
	if location.WebhookUrl == "" {
		location.WebhookUrl = defaults.WebhookUrl
	}

	location.WebhookAuthHeader = os.Getenv("WEBHOOK_AUTH_HEADER")
	if location.WebhookAuthHeader == "" {
		location.WebhookAuthHeader = defaults.WebhookAuthHeader
	}

	location.WebhookAuthToken = os.Getenv("WEBHOOK_AUTH_TOKEN")
	if location.WebhookAuthToken == "" {
		location.WebhookAuthToken = defaults.WebhookAuthToken
	}

	return location
}

func InitConfig() {
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}
	Config.HttpPort = httpPort

	powerOffTimeout := os.Getenv("POWER_OFF_TIMEOUT")
	if powerOffTimeout == "" {
		powerOffTimeout = "1m"
	}
	Config.PowerOffTimeout = powerOffTimeout

	locationIds := strings.Split(os.Getenv("LOCATION_IDS"), ",")

	Config.Locations = make(map[string]location)
	locationDefaults := getLocationConfigFromEnv("DEFAULT", location{})

	Config.Locations["DEFAULT"] = locationDefaults

	for _, locationId := range locationIds {
		if locationId == "DEFAULT" {
			return
		}
		Config.Locations[locationId] = getLocationConfigFromEnv(locationId, locationDefaults)
	}

	log.Println("generated global config from env variables") // @TODO: print the config without credentials
}
