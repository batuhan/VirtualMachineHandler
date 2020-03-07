package helpers

import (
	"encoding/json"
	"log"
)

type VmInfoDump struct {
	VirtualMachines []struct {
		Self struct {
			Value string
		}
	}
}

func GetVCenterIdFromJSON(output []byte, logger *log.Logger) string {
	var dump VmInfoDump
	err := json.Unmarshal(output, &dump)
	if err != nil {
		logger.Println(err.Error())
	}
	return dump.VirtualMachines[0].Self.Value
}
