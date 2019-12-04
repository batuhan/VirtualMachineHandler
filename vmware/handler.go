package vmware

import (
	"fmt"
	"github.com/nodgroup/VirtualMachineHandler/helpers"
	"net/http"
	"os"
	"strconv"
)

func Env(w http.ResponseWriter, req *http.Request) {
	body := helpers.GetBody(req.Body)

	out, err := execute(body.Identifier, "env")
	if err != nil {
		http.Error(w, err.Error()+"\n"+string(out), http.StatusBadRequest)
		return
	}

	_, _ = fmt.Fprintf(w, string(out))
}

func Create(w http.ResponseWriter, req *http.Request) {
	body := helpers.GetBody(req.Body)

	out, err := execute(body.Identifier, "vm.clone", "-vm="+body.Template, "-on=false",
		"-c="+strconv.Itoa(body.Cpu), "-m="+strconv.Itoa(body.Memory), body.TargetName)
	if err != nil {
		http.Error(w, err.Error()+"\n"+string(out), http.StatusBadRequest)
		return
	}

	out, err = execute(body.Identifier, "object.mv", "./vm/"+body.TargetName,
		os.Getenv(body.Identifier+"_TARGET_DIRECTORY"))
	if err != nil {
		http.Error(w, err.Error()+"\n"+string(out), http.StatusBadRequest)
		return
	}

	out, err = execute(body.Identifier, "vm.disk.change", "-vm="+body.TargetName, "-size="+body.DiskSize)
	if err != nil {
		http.Error(w, err.Error()+"\n"+string(out), http.StatusBadRequest)
		return
	}

	out, err = execute(body.Identifier, "vm.power", "-on=true", body.TargetName)
	if err != nil {
		http.Error(w, err.Error()+"\n"+string(out), http.StatusBadRequest)
		return
	}

	_, _ = fmt.Fprintf(w, "OK")
}
