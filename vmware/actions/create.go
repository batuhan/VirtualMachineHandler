package actions

import (
	"encoding/base64"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sethvargo/go-password/password"
	"gitlab.com/nod/bigcore/VirtualMachineHandler/helpers"
	"gitlab.com/nod/bigcore/VirtualMachineHandler/vmware"
	"gopkg.in/yaml.v2"
	"strconv"
	"strings"
	"time"
)

func Create(body helpers.Create, uuid uuid.UUID) {
	logger := helpers.CreateLogger(body.Identifier + " " + body.TargetName)

	currentDate := time.Now().UTC().Format(time.RFC3339)
	vmName := currentDate + "_" + body.TargetName
	vmName = helpers.ApplyCreateNameRegex(vmName)

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

	out, err := vmware.Execute(body.Identifier, true, logger, "vm.clone", "-vm="+body.Template,
		"-on=false", "-c="+strconv.Itoa(body.Cpu), "-m="+strconv.Itoa(body.Memory), vmName)
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

	out, err = vmware.Execute(body.Identifier, true, logger, "object.mv", "./vm/"+vmName,
		helpers.Config.DynamicConfigs[body.Identifier].TargetDirectory)
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

	out, err = vmware.Execute(body.Identifier, true, logger, "vm.disk.change", "-vm="+vmName,
		"-size="+body.DiskSize)
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
		out, err = vmware.Execute(body.Identifier, true, logger, "vm.change", "-vm="+vmName,
			"-e=guestinfo.metadata=\""+base64.StdEncoding.EncodeToString(metadataString)+"\"",
			"-e=guestinfo.metadata.encoding=base64")
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

	out, err = vmware.Execute(body.Identifier, true, logger, "vm.change", "-vm="+vmName,
		"-e=guestinfo.userdata=\""+base64.StdEncoding.EncodeToString(userData)+"\"",
		"-e=guestinfo.userdata.encoding=base64")
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

	out, err = vmware.Execute(body.Identifier, true, logger, "vm.power", "-on=true", vmName)
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
	out, err = vmware.GetVMInfoDump(body.Identifier, vmName, logger)
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
		VCenterId: helpers.ParseVMInfoDump(out, logger).VirtualMachines[0].Self.Value,
	}, logger)
}
