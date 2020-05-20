package actions

import (
	"VirtualMachineHandler/helpers"
	"VirtualMachineHandler/vmware"
	"github.com/google/uuid"
	"strconv"
	"time"
)

//func GetPowerState(body helpers.Update, logger *log.Logger, uuid uuid.UUID) *string {
//	out, err := vmware.GetVMInfoDump(body.LocationId, body.TargetName, logger)
//	if err != nil {
//		logger.Println(err.Error())
//		logger.Println(string(out))
//		go helpers.SendWebhook(helpers.Webhook{
//			Uuid:             uuid.String(),
//			Step:             "getPowerState",
//			Success:          false,
//			ErrorExplanation: err.Error() + "\n" + string(out),
//		}, logger)
//		return nil
//	}
//
//	dump := helpers.ParseVMInfoDump(out, logger)
//	return &dump.VirtualMachines[0].Summary.Runtime.PowerState
//}

func Update(body helpers.Update, uuid uuid.UUID) {
	logger := helpers.CreateLogger(body.LocationId + " " + body.TargetName)

	baseVMName := helpers.ApplyTargetNameRegex(body.TargetName)
	vmName, err := vmware.FindVM(body.LocationId, logger, baseVMName, uuid)
	if err != nil {
		return
	}

	err = vmware.PowerOffVM(body.LocationId, vmName, logger, uuid)
	if err != nil {
		return
	}

	if body.Cpu != 0 {
		out, err := vmware.Execute(body.LocationId, true, logger, "vm.change", "-vm="+vmName,
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
		out, err := vmware.Execute(body.LocationId, true, logger, "vm.change", "-vm="+vmName,
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
		out, err := vmware.Execute(body.LocationId, true, logger, "vm.disk.change", "-vm="+vmName,
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

	out, err := vmware.Execute(body.LocationId, true, logger, "vm.power", "-on=true", vmName)
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
