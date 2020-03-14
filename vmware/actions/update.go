package actions

import (
	"github.com/google/uuid"
	"gitlab.com/nod/bigcore/VirtualMachineHandler/helpers"
	"gitlab.com/nod/bigcore/VirtualMachineHandler/vmware"
	"log"
	"strconv"
	"time"
)

func GetPowerState(body helpers.Update, logger *log.Logger, uuid uuid.UUID) *string {
	out, err := vmware.GetVMInfoDump(body.Identifier, body.TargetName, logger)
	if err != nil {
		logger.Println(err.Error())
		logger.Println(string(out))
		go helpers.SendWebhook(helpers.Webhook{
			Uuid:             uuid.String(),
			Step:             "getPowerState",
			Success:          false,
			ErrorExplanation: err.Error() + "\n" + string(out),
		}, logger)
		return nil
	}

	dump := helpers.ParseVMInfoDump(out, logger)
	return &dump.VirtualMachines[0].Summary.Runtime.PowerState
}

func Update(body helpers.Update, uuid uuid.UUID) {
	logger := helpers.CreateLogger(body.Identifier + " " + body.TargetName)

	out, err := vmware.Execute(body.Identifier, true, logger, "vm.power", "-on=true", body.TargetName)
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

	loopCount := 0
	powerState := GetPowerState(body, logger, uuid)
	if powerState == nil {
		return
	}

	for *powerState != "poweredOff" {
		loopCount++
		if loopCount >= 10 {
			logger.Println("can't get status after " + strconv.Itoa(loopCount) + " tries, aborting")
			return
		}

		logger.Println("power state is: " + *powerState)
		powerState = GetPowerState(body, logger, uuid)
	}
	go helpers.SendWebhook(helpers.Webhook{
		Uuid:    uuid.String(),
		Step:    "powerOffVM",
		Success: true,
	}, logger)

	if body.Cpu != 0 {
		out, err := vmware.Execute(body.Identifier, true, logger, "vm.change", "-vm="+body.TargetName, "-c="+strconv.Itoa(body.Cpu))
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
		out, err := vmware.Execute(body.Identifier, true, logger, "vm.change", "-vm="+body.TargetName, "-m="+strconv.Itoa(body.Memory))
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
		out, err := vmware.Execute(body.Identifier, true, logger, "vm.disk.change", "-vm="+body.TargetName, "-size="+body.DiskSize)
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

	out, err = vmware.Execute(body.Identifier, true, logger, "vm.power", "-on=true", body.TargetName)
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
