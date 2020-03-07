package vmware

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/sethvargo/go-password/password"
	"gitlab.com/nod/bigcore/VirtualMachineHandler/helpers"
	"gopkg.in/yaml.v2"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func Env(w http.ResponseWriter, req *http.Request) {
	body := helpers.GetBody(req.Body)
	logger := helpers.CreateLogger(body.Identifier)

	out, err := execute(body.Identifier, true, logger, "env")
	if err != nil {
		logger.Println(err.Error())
		logger.Println(string(out))
		return
	}

	_, _ = fmt.Fprintf(w, string(out))
}

func Create(body helpers.Create, uuid uuid.UUID) {
	logger := helpers.CreateLogger(body.Identifier + " " + body.TargetName)

	pass, err := password.Generate(12, 2, 2, false, false)
	if err != nil {
		logger.Println(err.Error())
		go helpers.SendWebhook(helpers.Webhook{
			Uuid:             uuid.String(),
			Step:             "passwordGeneration",
			Success:          false,
			ErrorExplanation: "INTERNAL ERROR",
		}, logger)
		return
	}

	networkTemplate := helpers.CreateNetworkTemplate(body.Identifier, body.IpToAssign)

	template := helpers.GenerateBaseTemplate(body.SshKey, body.OnFirstBoot)
	var metadata *helpers.Metadata
	lowerTemplate := strings.ToLower(body.Template)
	isUbuntu := strings.Contains(lowerTemplate, "ubuntu")
	isDebian := strings.Contains(lowerTemplate, "debian")
	isCentos7 := strings.Contains(lowerTemplate, "centos-7")
	isCentos8 := strings.Contains(lowerTemplate, "centos-8")
	if isUbuntu {
		template, _ = helpers.AddSpecificParameters("ubuntu", template, pass, networkTemplate)
	} else if isDebian {
		template, metadata = helpers.AddSpecificParameters("debian", template, pass, networkTemplate)
	} else if isCentos7 {
		template, metadata = helpers.AddSpecificParameters("centos-7", template, pass, networkTemplate)
	} else if isCentos8 {
		template, metadata = helpers.AddSpecificParameters("centos-8", template, pass, networkTemplate)
	}

	userData, err := yaml.Marshal(template)
	if err != nil {
		logger.Println(err.Error())
		go helpers.SendWebhook(helpers.Webhook{
			Uuid:             uuid.String(),
			Step:             "templateGeneration",
			Success:          false,
			ErrorExplanation: "INTERNAL ERROR",
		}, logger)
		return
	}
	userData = append([]byte("#cloud-config\n"), userData...)

	var metadataString []byte
	if isCentos7 || isCentos8 || isDebian {
		metadataString, err = json.Marshal(metadata)
		if err != nil {
			logger.Println(err.Error())
			go helpers.SendWebhook(helpers.Webhook{
				Uuid:             uuid.String(),
				Step:             "templateGeneration",
				Success:          false,
				ErrorExplanation: "INTERNAL ERROR",
			}, logger)
			return
		}
	}

	go helpers.SendWebhook(helpers.Webhook{
		Uuid:    uuid.String(),
		Step:    "templateGeneration",
		Success: true,
	}, logger)

	out, err := execute(body.Identifier, true, logger, "vm.clone", "-vm="+body.Template, "-on=false",
		"-c="+strconv.Itoa(body.Cpu), "-m="+strconv.Itoa(body.Memory), body.TargetName)
	if err != nil {
		logger.Println(err.Error())
		logger.Println(string(out))
		go helpers.SendWebhook(helpers.Webhook{
			Uuid:             uuid.String(),
			Step:             "createVM",
			Success:          false,
			ErrorExplanation: err.Error() + "\n" + string(out),
		}, logger)
		return
	}
	go helpers.SendWebhook(helpers.Webhook{
		Uuid:    uuid.String(),
		Step:    "createVM",
		Success: true,
	}, logger)

	out, err = execute(body.Identifier, true, logger, "object.mv", "./vm/"+body.TargetName,
		os.Getenv(body.Identifier+"_TARGET_DIRECTORY"))
	if err != nil {
		logger.Println(err.Error())
		logger.Println(string(out))
		go helpers.SendWebhook(helpers.Webhook{
			Uuid:             uuid.String(),
			Step:             "moveVMToTargetDirectory",
			Success:          false,
			ErrorExplanation: err.Error() + "\n" + string(out),
		}, logger)
		return
	}
	go helpers.SendWebhook(helpers.Webhook{
		Uuid:    uuid.String(),
		Step:    "moveVMToTargetDirectory",
		Success: true,
	}, logger)

	out, err = execute(body.Identifier, true, logger, "vm.disk.change", "-vm="+body.TargetName, "-size="+body.DiskSize)
	if err != nil {
		logger.Println(err.Error())
		logger.Println(string(out))
		go helpers.SendWebhook(helpers.Webhook{
			Uuid:             uuid.String(),
			Step:             "changeVMDiskSize",
			Success:          false,
			ErrorExplanation: err.Error() + "\n" + string(out),
		}, logger)
		return
	}
	go helpers.SendWebhook(helpers.Webhook{
		Uuid:    uuid.String(),
		Step:    "changeVMDiskSize",
		Success: true,
	}, logger)

	if isCentos7 || isCentos8 || isDebian {
		out, err = execute(body.Identifier, true, logger, "vm.change", "-vm="+body.TargetName,
			"-e=guestinfo.metadata=\""+base64.StdEncoding.EncodeToString(metadataString)+"\"", "-e=guestinfo.metadata.encoding=base64")
		if err != nil {
			logger.Println(err.Error())
			logger.Println(string(out))
			go helpers.SendWebhook(helpers.Webhook{
				Uuid:             uuid.String(),
				Step:             "addCloudInitTemplate",
				Success:          false,
				ErrorExplanation: err.Error() + "\n" + string(out),
			}, logger)
			return
		}
	}

	out, err = execute(body.Identifier, true, logger, "vm.change", "-vm="+body.TargetName,
		"-e=guestinfo.userdata=\""+base64.StdEncoding.EncodeToString(userData)+"\"", "-e=guestinfo.userdata.encoding=base64")
	if err != nil {
		logger.Println(err.Error())
		logger.Println(string(out))
		go helpers.SendWebhook(helpers.Webhook{
			Uuid:             uuid.String(),
			Step:             "addCloudInitTemplate",
			Success:          false,
			ErrorExplanation: err.Error() + "\n" + string(out),
		}, logger)
		return
	}
	go helpers.SendWebhook(helpers.Webhook{
		Uuid:    uuid.String(),
		Step:    "addCloudInitTemplate",
		Success: true,
	}, logger)

	out, err = execute(body.Identifier, true, logger, "vm.power", "-on=true", body.TargetName)
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

	// get vCenter ID
	// /VirtualMachines/0/Self/Value
	out, err = execute(body.Identifier, true, logger, "vm.info", "-json", body.TargetName)
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
		Uuid:      uuid.String(),
		Step:      "powerOnVM",
		Success:   true,
		Password:  pass,
		VCenterId: helpers.GetVCenterIdFromJSON(out, logger),
	}, logger)
}

func Delete(identifier string, targetName string, uuid uuid.UUID) error {
	logger := helpers.CreateLogger(identifier + " " + targetName)

	out, err := execute(identifier, true, logger, "vm.destroy", targetName)
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

func Recreate(body helpers.Create, uuid uuid.UUID) {
	err := Delete(body.Identifier, body.TargetName, uuid)
	if err != nil {
		return
	}
	Create(body, uuid)
}

func Update(body helpers.Update, uuid uuid.UUID) {
	logger := helpers.CreateLogger(body.Identifier + " " + body.TargetName)

	out, err := execute(body.Identifier, true, logger, "vm.power", "-s=true", body.TargetName)
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
	go helpers.SendWebhook(helpers.Webhook{
		Uuid:    uuid.String(),
		Step:    "powerOffVM",
		Success: true,
	}, logger)

	if body.Cpu != 0 {
		out, err := execute(body.Identifier, true, logger, "vm.change", "-vm="+body.TargetName, "-c="+strconv.Itoa(body.Cpu))
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
		out, err := execute(body.Identifier, true, logger, "vm.change", "-vm="+body.TargetName, "-m="+strconv.Itoa(body.Memory))
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
		out, err := execute(body.Identifier, true, logger, "vm.disk.change", "-vm="+body.TargetName, "-size="+body.DiskSize)
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

	out, err = execute(body.Identifier, true, logger, "vm.power", "-on=true", body.TargetName)
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

	out, err := execute(body.Identifier, true, logger, "vm.power", "-"+nextState+"=true", body.TargetName)
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
