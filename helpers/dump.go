package helpers

import (
	"encoding/json"
	"log"
)

type VMInfoDump struct {
	VirtualMachines []struct {
		Summary struct {
			Runtime struct {
				PowerState string
			}
		}
		Self struct {
			Value string
		}
		Config struct {
			Hardware struct {
				Device []struct {
					MacAddress string
				}
			}
		}
	}
}

func ParseVMInfoDump(output []byte, logger *log.Logger) VMInfoDump {
	var dump VMInfoDump
	err := json.Unmarshal(output, &dump)
	if err != nil {
		logger.Println(err.Error())
	}
	return dump
}

func GetMacAddress(devices []struct{ MacAddress string }) string {
	for _, device := range devices {
		macAddress := device.MacAddress
		if macAddress != "" {
			return macAddress
		}
	}
	return ""
}
