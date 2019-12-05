package vmware

import (
	"encoding/base64"
	"fmt"
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

	out, err := execute(body.Identifier, "env")
	if err != nil {
		log.Println(err.Error())
		log.Println(string(out))
		return
	}

	_, _ = fmt.Fprintf(w, string(out))
}

func Create(w http.ResponseWriter, req *http.Request) {
	body := helpers.GetBody(req.Body)

	networkTemplateObject := helpers.CreateNetworkTemplate(body.Identifier, body.IpToAssign)
	networkTemplate, err := yaml.Marshal(networkTemplateObject)
	if err != nil {
		log.Println(err.Error())
		return
	}

	template := helpers.GenerateBaseTemplate(body.SshKey)
	if strings.Contains(strings.ToLower(body.Template), "ubuntu") {
		template, err = helpers.AddUbuntuSpecificParameters(template, networkTemplate)
	}
	if err != nil {
		log.Println(err.Error())
		return
	}

	userData, err := yaml.Marshal(template)
	if err != nil {
		log.Println(err.Error())
		return
	}
	userData = append([]byte("#cloud-config\n"), userData...)

	out, err := execute(body.Identifier, "vm.clone", "-vm="+body.Template, "-on=false",
		"-c="+strconv.Itoa(body.Cpu), "-m="+strconv.Itoa(body.Memory), body.TargetName)
	if err != nil {
		log.Println(err.Error())
		log.Println(string(out))
		return
	}

	out, err = execute(body.Identifier, "object.mv", "./vm/"+body.TargetName,
		os.Getenv(body.Identifier+"_TARGET_DIRECTORY"))
	if err != nil {
		log.Println(err.Error())
		log.Println(string(out))
		return
	}

	out, err = execute(body.Identifier, "vm.disk.change", "-vm="+body.TargetName, "-size="+body.DiskSize)
	if err != nil {
		log.Println(err.Error())
		log.Println(string(out))
		return
	}

	out, err = execute(body.Identifier, "vm.change", "-vm="+body.TargetName,
		"-e=guestinfo.userdata=\""+base64.StdEncoding.EncodeToString(userData)+"\"", "-e=guestinfo.userdata.encoding=base64")
	if err != nil {
		log.Println(err.Error())
		log.Println(string(out))
		return
	}

	out, err = execute(body.Identifier, "vm.power", "-on=true", body.TargetName)
	if err != nil {
		log.Println(err.Error())
		log.Println(string(out))
		return
	}

	_, _ = fmt.Fprintf(w, "OK")
}
