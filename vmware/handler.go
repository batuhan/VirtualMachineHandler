package vmware

import (
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"gitlab.com/nod/bigcore/VirtualMachineHandler/helpers"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	networkTemplateObject := helpers.CreateNetworkTemplate(body.Identifier, body.IpToAssign)
	networkTemplate, err := yaml.Marshal(networkTemplateObject)
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

	template := helpers.GenerateBaseTemplate(body.SshKey)
	if strings.Contains(strings.ToLower(body.Template), "ubuntu") {
		template, err = helpers.AddUbuntuSpecificParameters(template, networkTemplate)
	}
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

	helpers.SendWebhook(helpers.Webhook{
		Uuid:             uuid.String(),
		Step:             "templateGeneration",
		Success:          true,
		ErrorExplanation: "",
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
		Uuid:             uuid.String(),
		Step:             "createVM",
		Success:          true,
		ErrorExplanation: "",
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
		Uuid:             uuid.String(),
		Step:             "moveVMToTargetDirectory",
		Success:          true,
		ErrorExplanation: "",
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
		Uuid:             uuid.String(),
		Step:             "changeVMDiskSize",
		Success:          true,
		ErrorExplanation: "",
	})

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
		Uuid:             uuid.String(),
		Step:             "addCloudInitTemplate",
		Success:          true,
		ErrorExplanation: "",
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
		Uuid:             uuid.String(),
		Step:             "powerOnVM",
		Success:          true,
		ErrorExplanation: "",
	})
}
