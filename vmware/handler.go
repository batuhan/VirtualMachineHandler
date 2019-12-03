package vmware

import (
	"fmt"
	"github.com/nodgroup/VirtualMachineHandler/helpers"
	"net/http"
)

func Env(w http.ResponseWriter, req *http.Request) {
	body := helpers.GetBody(req.Body)
	out, err := execute(body.Identifier, "env")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	_, _ = fmt.Fprintf(w, string(out))
}
