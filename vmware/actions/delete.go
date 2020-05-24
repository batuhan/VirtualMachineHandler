package actions

import (
	"VirtualMachineHandler/helpers"
	"VirtualMachineHandler/vmware"
	"github.com/google/uuid"
)

func Delete(identifier string, targetName string, uuid uuid.UUID) error {
	logger := helpers.CreateLogger(identifier + " " + targetName)

	baseVMName := helpers.ApplyTargetNameRegex(targetName)
	vmName, err := vmware.FindVM(identifier, logger, baseVMName, uuid)
	if err != nil {
		return err
	}

	err = vmware.PowerOffVM(identifier, vmName, logger, uuid)
	if err != nil {
		return err
	}

	out, err := vmware.Execute(identifier, true, logger, "object.mv", vmName,
		helpers.Config.Locations[identifier].DeleteDirectory)
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
