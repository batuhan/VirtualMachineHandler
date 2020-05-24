package vmware

import (
	"VirtualMachineHandler/helpers"
	"log"
	"os/exec"
)

func Execute(baseEnv string, doLog bool, logger *log.Logger, arg ...string) ([]byte, error) {
	cmd := exec.Command("govc", arg...)
	if doLog {
		logger.Println(cmd.Args)
	}

	location := helpers.Config.Locations[baseEnv]
	cmd.Env = []string{
		"GOVC_INSECURE=" + location.GovcInsecure,
		"GOVC_URL=" + location.GovcURL,
		"GOVC_USERNAME=" + location.GovcUsername,
		"GOVC_PASSWORD=" + location.GovcPassword,
		"GOVC_DATACENTER=" + location.GovcDatacenter,
		"GOVC_DATASTORE=" + location.GovcDatastore,
		"GOVC_RESOURCE_POOL=" + location.GovcResourcePool,
	}
	out, err := cmd.CombinedOutput()
	// safeOutput := strings.Replace(string(out), location.GovcPassword, "[HIDDEN]", -1)
	return out, err
}

func GetVMInfoDump(identifier string, targetName string, logger *log.Logger) ([]byte, error) {
	out, err := Execute(identifier, true, logger, "vm.info", "-json", targetName)
	return out, err
}
