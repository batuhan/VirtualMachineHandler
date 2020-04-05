package vmware

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gitlab.com/nod/bigcore/VirtualMachineHandler/helpers"
	"log"
	"strings"
)

func FindVM(id string, logger *log.Logger, baseVMName string, uuid uuid.UUID) (string, error) {
	out, err := Execute(id, true, logger, "find", helpers.Config.Locations[id].TargetDirectory,
		"-type", "m", "-name", "*"+baseVMName)
	if err != nil {
		logger.Println(err.Error())
		logger.Println(string(out))
		go helpers.SendWebhook(helpers.Webhook{
			Uuid:             uuid.String(),
			Step:             "findVM",
			Success:          false,
			ErrorExplanation: err.Error() + "\n" + string(out),
		}, logger)
		return "", err
	}
	go helpers.SendWebhook(helpers.Webhook{
		Uuid:    uuid.String(),
		Step:    "findVM",
		Success: true,
	}, logger)
	vmName := string(out)
	split := strings.Split(vmName, "/")
	vmName = strings.TrimSpace(split[len(split)-1])

	return vmName, nil
}

func PowerOffVM(identifier string, vmName string, logger *log.Logger, uuid uuid.UUID) error {
	out, err := Execute(identifier, true, logger, "vm.power", "-off=true", vmName)
	if err != nil {
		logger.Println(err.Error())
		logger.Println(string(out))
		// disable webhook since vm power off is a pre-condition and a already closed vm will throw an error
		//go helpers.SendWebhook(helpers.Webhook{
		//	Uuid:             uuid.String(),
		//	Step:             "powerOffVM",
		//	Success:          false,
		//	ErrorExplanation: err.Error() + "\n" + string(out),
		//})
	}

	out, err = Execute(identifier, true, logger, "object.collect", "-s",
		fmt.Sprintf("-wait=%s", helpers.Config.PowerOffTimeout),
		fmt.Sprintf("%s/%s", helpers.Config.Locations[identifier].TargetDirectory, vmName),
		"-runtime.powerState", "poweredOff")
	if err != nil || len(out) == 0 {
		var errorExplanation string
		var powerOffError error
		if err != nil {
			logger.Println(err.Error())
			errorExplanation = err.Error() + "\n" + string(out)
			powerOffError = err
		} else {
			errorExplanation = "can't verify power status after timeout"
			powerOffError = errors.New(errorExplanation)
		}
		logger.Println(string(out))
		go helpers.SendWebhook(helpers.Webhook{
			Uuid:             uuid.String(),
			Step:             "powerOffVM",
			Success:          false,
			ErrorExplanation: errorExplanation,
		}, logger)
		return powerOffError
	}

	go helpers.SendWebhook(helpers.Webhook{
		Uuid:    uuid.String(),
		Step:    "powerOffVM",
		Success: true,
	}, logger)

	return nil
}
