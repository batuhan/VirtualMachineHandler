package actions

import (
	"fmt"
	"gitlab.com/nod/bigcore/VirtualMachineHandler/helpers"
	"gitlab.com/nod/bigcore/VirtualMachineHandler/vmware"
	"net/http"
	"strings"
)

func Env(w http.ResponseWriter, req *http.Request) {
	body := helpers.GetBody(req.Body)

	logger := helpers.CreateLogger(body.LocationId)

	out, err := vmware.Execute(body.LocationId, true, logger, "env")
	safeOutput := strings.Replace(string(out), helpers.Config.Locations[body.LocationId].GovcPassword, "[HIDDEN]", -1)

	if err != nil {
		logger.Println(err.Error())
		logger.Println(safeOutput)
		return
	}

	_, _ = fmt.Fprintf(w, safeOutput)
}
