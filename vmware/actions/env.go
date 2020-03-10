package actions

import (
	"fmt"
	"gitlab.com/nod/bigcore/VirtualMachineHandler/helpers"
	"gitlab.com/nod/bigcore/VirtualMachineHandler/vmware"
	"net/http"
)

func Env(w http.ResponseWriter, req *http.Request) {
	body := helpers.GetBody(req.Body)
	logger := helpers.CreateLogger(body.Identifier)

	out, err := vmware.Execute(body.Identifier, true, logger, "env")
	if err != nil {
		logger.Println(err.Error())
		logger.Println(string(out))
		return
	}

	_, _ = fmt.Fprintf(w, string(out))
}
