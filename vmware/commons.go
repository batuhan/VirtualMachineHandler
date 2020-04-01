package vmware

import (
	"github.com/google/uuid"
	"gitlab.com/nod/bigcore/VirtualMachineHandler/helpers"
	"log"
	"strings"
)

func FindVM(id string, logger *log.Logger, baseVMName string, uuid uuid.UUID) (string, error) {
	out, err := Execute(id, true, logger, "find", ".", "-type", "m", "-name", "*"+baseVMName)
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
