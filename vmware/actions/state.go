package actions

import (
	"github.com/google/uuid"
	"gitlab.com/nod/bigcore/VirtualMachineHandler/helpers"
	"gitlab.com/nod/bigcore/VirtualMachineHandler/vmware"
)

func State(body helpers.State, uuid uuid.UUID) {
	logger := helpers.CreateLogger(body.Identifier + " " + body.TargetName)

	nextState := ""

	if body.Action == "on" {
		nextState = "on"
	} else if body.Action == "off" {
		nextState = "off"
	} else if body.Action == "suspend" {
		nextState = "suspend"
	} else if body.Action == "reset" {
		nextState = "reset"
	} else if body.Action == "shutdown" {
		nextState = "s"
	} else if body.Action == "reboot" {
		nextState = "r"
	} else {
		go helpers.SendWebhook(helpers.Webhook{
			Uuid:             uuid.String(),
			Step:             "VMStateChange",
			Success:          false,
			ErrorExplanation: "need valid action value",
		}, logger)
		return
	}

	out, err := vmware.Execute(body.Identifier, true, logger, "vm.power", "-"+nextState+"=true", body.TargetName)
	if err != nil {
		logger.Println(err.Error())
		logger.Println(string(out))
		go helpers.SendWebhook(helpers.Webhook{
			Uuid:             uuid.String(),
			Step:             "VMStateChange",
			Success:          false,
			ErrorExplanation: err.Error() + "\n" + string(out),
		}, logger)
		return
	}
	go helpers.SendWebhook(helpers.Webhook{
		Uuid:    uuid.String(),
		Step:    "VMStateChange",
		Success: true,
	}, logger)
}