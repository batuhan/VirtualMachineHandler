package actions

import (
	"github.com/google/uuid"
	"gitlab.com/nod/bigcore/VirtualMachineHandler/helpers"
	"gitlab.com/nod/bigcore/VirtualMachineHandler/vmware"
	"strconv"
	"time"
)

func Update(body helpers.Update, uuid uuid.UUID) {
	logger := helpers.CreateLogger(body.Identifier + " " + body.TargetName)

	baseVMName := helpers.ApplyTargetNameRegex(body.TargetName)
	vmName, err := vmware.FindVM(body.Identifier, logger, baseVMName, uuid)
	if err != nil {
		return
	}

	err = vmware.PowerOffVM(body.Identifier, vmName, logger, uuid)
	if err != nil {
		return
	}

	if body.Cpu != 0 {
		out, err := vmware.Execute(body.Identifier, true, logger, "vm.change", "-vm="+vmName,
			"-c="+strconv.Itoa(body.Cpu))
		if err != nil {
			logger.Println(err.Error())
			logger.Println(string(out))
			go helpers.SendWebhook(helpers.Webhook{
				Uuid:             uuid.String(),
				Step:             "updateCPU",
				Success:          false,
				ErrorExplanation: err.Error() + "\n" + string(out),
			}, logger)
			return
		}
		go helpers.SendWebhook(helpers.Webhook{
			Uuid:    uuid.String(),
			Step:    "updateCPU",
			Success: true,
		}, logger)
	}

	if body.Memory != 0 {
		out, err := vmware.Execute(body.Identifier, true, logger, "vm.change", "-vm="+vmName,
			"-m="+strconv.Itoa(body.Memory))
		if err != nil {
			logger.Println(err.Error())
			logger.Println(string(out))
			go helpers.SendWebhook(helpers.Webhook{
				Uuid:             uuid.String(),
				Step:             "updateMemory",
				Success:          false,
				ErrorExplanation: err.Error() + "\n" + string(out),
			}, logger)
			return
		}
		go helpers.SendWebhook(helpers.Webhook{
			Uuid:    uuid.String(),
			Step:    "updateMemory",
			Success: true,
		}, logger)
	}

	if body.DiskSize != "" {
		out, err := vmware.Execute(body.Identifier, true, logger, "vm.disk.change", "-vm="+vmName,
			"-size="+body.DiskSize)
		time.Sleep(5 * time.Second)
		if err != nil {
			logger.Println(err.Error())
			logger.Println(string(out))
			// disable error webhook for disk size since shrinking disk will always result in an error
			//go helpers.SendWebhook(helpers.Webhook{
			//	Uuid:             uuid.String(),
			//	Step:             "updateDiskSize",
			//	Success:          false,
			//	ErrorExplanation: err.Error() + "\n" + string(out),
			//})
		}
		go helpers.SendWebhook(helpers.Webhook{
			Uuid:    uuid.String(),
			Step:    "updateDiskSize",
			Success: true,
		}, logger)
	}

	out, err := vmware.Execute(body.Identifier, true, logger, "vm.power", "-on=true", vmName)
	if err != nil {
		logger.Println(err.Error())
		logger.Println(string(out))
		go helpers.SendWebhook(helpers.Webhook{
			Uuid:             uuid.String(),
			Step:             "powerOnVM",
			Success:          false,
			ErrorExplanation: err.Error() + "\n" + string(out),
		}, logger)
		return
	}
	go helpers.SendWebhook(helpers.Webhook{
		Uuid:    uuid.String(),
		Step:    "powerOnVM",
		Success: true,
	}, logger)
	go helpers.SendWebhook(helpers.Webhook{
		Uuid:    uuid.String(),
		Step:    "updateFinished",
		Success: true,
	}, logger)
}
