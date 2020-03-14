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
