package actions

import (
	"github.com/google/uuid"
	"gitlab.com/nod/bigcore/VirtualMachineHandler/helpers"
	"gitlab.com/nod/bigcore/VirtualMachineHandler/vmware"
)

func Delete(identifier string, targetName string, uuid uuid.UUID) error {
	logger := helpers.CreateLogger(identifier + " " + targetName)

	out, err := vmware.Execute(identifier, true, logger, "vm.destroy", targetName)
	if err != nil {
		logger.Println(err.Error())
		logger.Println(string(out))
		go helpers.SendWebhook(helpers.Webhook{
			Uuid:             uuid.String(),
			Step:             "deleteVM",
			Success:          false,
			ErrorExplanation: err.Error() + "\n" + string(out),
		}, logger)
		return err
	}
	go helpers.SendWebhook(helpers.Webhook{
		Uuid:    uuid.String(),
		Step:    "deleteVM",
		Success: true,
	}, logger)
	return nil
}
