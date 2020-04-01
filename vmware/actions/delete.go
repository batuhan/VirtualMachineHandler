package actions

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gitlab.com/nod/bigcore/VirtualMachineHandler/helpers"
	"gitlab.com/nod/bigcore/VirtualMachineHandler/vmware"
)

func Delete(identifier string, targetName string, uuid uuid.UUID) error {
	logger := helpers.CreateLogger(identifier + " " + targetName)

	baseVMName := helpers.ApplyTargetNameRegex(targetName)
	vmName, err := vmware.FindVM(identifier, logger, baseVMName, uuid)
	if err != nil {
		return err
	}

	out, err := vmware.Execute(identifier, true, logger, "vm.power", "-off=true", vmName)
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

	out, err = vmware.Execute(identifier, true, logger, "object.collect", "-s",
		fmt.Sprintf("-wait=%s", helpers.Config.PowerOffTimeout),
		fmt.Sprintf("%s/%s", helpers.Config.DynamicConfigs[identifier].TargetDirectory, vmName),
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

	out, err = vmware.Execute(identifier, true, logger, "object.mv", vmName,
		helpers.Config.DynamicConfigs[identifier].DeleteDirectory)
	if err != nil {
		logger.Println(err.Error())
		logger.Println(string(out))
		go helpers.SendWebhook(helpers.Webhook{
			Uuid:             uuid.String(),
			Step:             "archiveVM",
			Success:          false,
			ErrorExplanation: err.Error() + "\n" + string(out),
		}, logger)
		return err
	}
	go helpers.SendWebhook(helpers.Webhook{
		Uuid:    uuid.String(),
		Step:    "archiveVM",
		Success: true,
	}, logger)
	return nil
}
