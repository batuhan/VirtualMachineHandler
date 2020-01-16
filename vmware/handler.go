package vmware

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/sethvargo/go-password/password"
	"gitlab.com/nod/bigcore/VirtualMachineHandler/helpers"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func Env(w http.ResponseWriter, req *http.Request) {
	body := helpers.GetBody(req.Body)

	out, err := execute(body.Identifier, true, "env")
	if err != nil {
		log.Println(err.Error())
		log.Println(string(out))
		return
	}

	_, _ = fmt.Fprintf(w, string(out))
}

func Create(body helpers.Body, uuid uuid.UUID) {
	pass, err := password.Generate(12, 2, 2, false, false)
	if err != nil {
		log.Println(err.Error())
		helpers.SendWebhook(helpers.Webhook{
			Uuid:             uuid.String(),
			Step:             "passwordGeneration",
			Success:          false,
			ErrorExplanation: "INTERNAL ERROR",
		})
		return
	}

	networkTemplate := helpers.CreateNetworkTemplate(body.Identifier, body.IpToAssign)

	template := helpers.GenerateBaseTemplate(body.SshKey)
	var metadata *helpers.Metadata
	lowerTemplate := strings.ToLower(body.Template)
	isUbuntu := strings.Contains(lowerTemplate, "ubuntu")
	isCentos := strings.Contains(lowerTemplate, "centos")
	if isUbuntu {
		template, _ = helpers.AddSpecificParameters("ubuntu", template, pass, networkTemplate)
	} else if isCentos {
		template, metadata = helpers.AddSpecificParameters("centos", template, pass, networkTemplate)
	}

	userData, err := yaml.Marshal(template)
	if err != nil {
		log.Println(err.Error())
		helpers.SendWebhook(helpers.Webhook{
			Uuid:             uuid.String(),
			Step:             "templateGeneration",
			Success:          false,
			ErrorExplanation: "INTERNAL ERROR",
		})
		return
	}
	userData = append([]byte("#cloud-config\n"), userData...)

	var metadataString []byte
	if isCentos {
		metadataString, err = json.Marshal(metadata)
		if err != nil {
			log.Println(err.Error())
			helpers.SendWebhook(helpers.Webhook{
				Uuid:             uuid.String(),
				Step:             "templateGeneration",
				Success:          false,
				ErrorExplanation: "INTERNAL ERROR",
			})
			return
		}
	}

	helpers.SendWebhook(helpers.Webhook{
		Uuid:    uuid.String(),
		Step:    "templateGeneration",
		Success: true,
	})

	out, err := execute(body.Identifier, true, "vm.clone", "-vm="+body.Template, "-on=false",
		"-c="+strconv.Itoa(body.Cpu), "-m="+strconv.Itoa(body.Memory), body.TargetName)
	if err != nil {
		log.Println(err.Error())
		log.Println(string(out))
		helpers.SendWebhook(helpers.Webhook{
			Uuid:             uuid.String(),
			Step:             "createVM",
			Success:          false,
			ErrorExplanation: err.Error() + "\n" + string(out),
		})
		return
	}
	helpers.SendWebhook(helpers.Webhook{
		Uuid:    uuid.String(),
		Step:    "createVM",
		Success: true,
	})

	out, err = execute(body.Identifier, true, "object.mv", "./vm/"+body.TargetName,
		os.Getenv(body.Identifier+"_TARGET_DIRECTORY"))
	if err != nil {
		log.Println(err.Error())
		log.Println(string(out))
		helpers.SendWebhook(helpers.Webhook{
			Uuid:             uuid.String(),
			Step:             "moveVMToTargetDirectory",
			Success:          false,
			ErrorExplanation: err.Error() + "\n" + string(out),
		})
		return
	}
	helpers.SendWebhook(helpers.Webhook{
		Uuid:    uuid.String(),
		Step:    "moveVMToTargetDirectory",
		Success: true,
	})

	out, err = execute(body.Identifier, true, "vm.disk.change", "-vm="+body.TargetName, "-size="+body.DiskSize)
	if err != nil {
		log.Println(err.Error())
		log.Println(string(out))
		helpers.SendWebhook(helpers.Webhook{
			Uuid:             uuid.String(),
			Step:             "changeVMDiskSize",
			Success:          false,
			ErrorExplanation: err.Error() + "\n" + string(out),
		})
		return
	}
	helpers.SendWebhook(helpers.Webhook{
		Uuid:    uuid.String(),
		Step:    "changeVMDiskSize",
		Success: true,
	})

	if isCentos {
		out, err = execute(body.Identifier, false, "vm.change", "-vm="+body.TargetName,
			"-e=guestinfo.metadata=\""+base64.StdEncoding.EncodeToString(metadataString)+"\"", "-e=guestinfo.metadata.encoding=base64")
		if err != nil {
			log.Println(err.Error())
			log.Println(string(out))
			helpers.SendWebhook(helpers.Webhook{
				Uuid:             uuid.String(),
				Step:             "addCloudInitTemplate",
				Success:          false,
				ErrorExplanation: err.Error() + "\n" + string(out),
			})
			return
		}
	}

	out, err = execute(body.Identifier, false, "vm.change", "-vm="+body.TargetName,
		"-e=guestinfo.userdata=\""+base64.StdEncoding.EncodeToString(userData)+"\"", "-e=guestinfo.userdata.encoding=base64")
	if err != nil {
		log.Println(err.Error())
		log.Println(string(out))
		helpers.SendWebhook(helpers.Webhook{
			Uuid:             uuid.String(),
			Step:             "addCloudInitTemplate",
			Success:          false,
			ErrorExplanation: err.Error() + "\n" + string(out),
		})
		return
	}
	helpers.SendWebhook(helpers.Webhook{
		Uuid:    uuid.String(),
		Step:    "addCloudInitTemplate",
		Success: true,
	})

	out, err = execute(body.Identifier, true, "vm.power", "-on=true", body.TargetName)
	if err != nil {
		log.Println(err.Error())
		log.Println(string(out))
		helpers.SendWebhook(helpers.Webhook{
			Uuid:             uuid.String(),
			Step:             "powerOnVM",
			Success:          false,
			ErrorExplanation: err.Error() + "\n" + string(out),
		})
		return
	}
	helpers.SendWebhook(helpers.Webhook{
		Uuid:     uuid.String(),
		Step:     "powerOnVM",
		Success:  true,
		Password: pass,
	})
}

func Delete(body helpers.Body, uuid uuid.UUID) error {
	out, err := execute(body.Identifier, true, "vm.destroy", body.TargetName)
	if err != nil {
		log.Println(err.Error())
		log.Println(string(out))
		helpers.SendWebhook(helpers.Webhook{
			Uuid:             uuid.String(),
			Step:             "deleteVM",
			Success:          false,
			ErrorExplanation: err.Error() + "\n" + string(out),
		})
		return err
	}
	helpers.SendWebhook(helpers.Webhook{
		Uuid:    uuid.String(),
		Step:    "deleteVM",
		Success: true,
	})
	return nil
}

func Recreate(body helpers.Body, uuid uuid.UUID) {
	err := Delete(body, uuid)
	if err != nil {
		return
	}
	Create(body, uuid)
}

func Update(body helpers.Body, uuid uuid.UUID) {
	out, err := execute(body.Identifier, true, "vm.power", "-s=true", body.TargetName)
	if err != nil {
		log.Println(err.Error())
		log.Println(string(out))
		// disable webhook since vm power off is a pre-condition and a already closed vm will throw an error
		//helpers.SendWebhook(helpers.Webhook{
		//	Uuid:             uuid.String(),
		//	Step:             "powerOffVM",
		//	Success:          false,
		//	ErrorExplanation: err.Error() + "\n" + string(out),
		//})
	}
	helpers.SendWebhook(helpers.Webhook{
		Uuid:    uuid.String(),
		Step:    "powerOffVM",
		Success: true,
	})

	if body.Cpu != 0 {
		out, err := execute(body.Identifier, true, "vm.change", "-vm="+body.TargetName, "-c="+strconv.Itoa(body.Cpu))
		if err != nil {
			log.Println(err.Error())
			log.Println(string(out))
			helpers.SendWebhook(helpers.Webhook{
				Uuid:             uuid.String(),
				Step:             "updateCPU",
				Success:          false,
				ErrorExplanation: err.Error() + "\n" + string(out),
			})
			return
		}
		helpers.SendWebhook(helpers.Webhook{
			Uuid:    uuid.String(),
			Step:    "updateCPU",
			Success: true,
		})
	}

	if body.Memory != 0 {
		out, err := execute(body.Identifier, true, "vm.change", "-vm="+body.TargetName, "-m="+strconv.Itoa(body.Memory))
		if err != nil {
			log.Println(err.Error())
			log.Println(string(out))
			helpers.SendWebhook(helpers.Webhook{
				Uuid:             uuid.String(),
				Step:             "updateMemory",
				Success:          false,
				ErrorExplanation: err.Error() + "\n" + string(out),
			})
			return
		}
		helpers.SendWebhook(helpers.Webhook{
			Uuid:    uuid.String(),
			Step:    "updateMemory",
			Success: true,
		})
	}

	if body.DiskSize != "" {
		out, err := execute(body.Identifier, true, "vm.disk.change", "-vm="+body.TargetName, "-size="+body.DiskSize)
		if err != nil {
			log.Println(err.Error())
			log.Println(string(out))
			// disable error webhook for disk size since shrinking disk will always result in an error
			//helpers.SendWebhook(helpers.Webhook{
			//	Uuid:             uuid.String(),
			//	Step:             "updateDiskSize",
			//	Success:          false,
			//	ErrorExplanation: err.Error() + "\n" + string(out),
			//})
		}
		helpers.SendWebhook(helpers.Webhook{
			Uuid:    uuid.String(),
			Step:    "updateDiskSize",
			Success: true,
		})
	}

	time.Sleep(time.Second)
	out, err = execute(body.Identifier, true, "vm.power", "-on=true", body.TargetName)
	if err != nil {
		log.Println(err.Error())
		log.Println(string(out))
		helpers.SendWebhook(helpers.Webhook{
			Uuid:             uuid.String(),
			Step:             "powerOnVM",
			Success:          false,
			ErrorExplanation: err.Error() + "\n" + string(out),
		})
		return
	}
	helpers.SendWebhook(helpers.Webhook{
		Uuid:    uuid.String(),
		Step:    "powerOnVM",
		Success: true,
	})
}

func State(body helpers.Body, uuid uuid.UUID) {
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
		helpers.SendWebhook(helpers.Webhook{
			Uuid:             uuid.String(),
			Step:             "VMStateChange",
			Success:          false,
			ErrorExplanation: "need valid action value",
		})
		return
	}

	out, err := execute(body.Identifier, true, "vm.power", "-"+nextState+"=true", body.TargetName)
	if err != nil {
		log.Println(err.Error())
		log.Println(string(out))
		helpers.SendWebhook(helpers.Webhook{
			Uuid:             uuid.String(),
			Step:             "VMStateChange",
			Success:          false,
			ErrorExplanation: err.Error() + "\n" + string(out),
		})
		return
	}
	helpers.SendWebhook(helpers.Webhook{
		Uuid:    uuid.String(),
		Step:    "VMStateChange",
		Success: true,
	})
}
