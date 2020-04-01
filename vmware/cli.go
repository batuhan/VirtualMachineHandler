package vmware

import (
	"gitlab.com/nod/bigcore/VirtualMachineHandler/helpers"
	"log"
	"os/exec"
)

func Execute(baseEnv string, doLog bool, logger *log.Logger, arg ...string) ([]byte, error) {
	cmd := exec.Command("govc", arg...)
	if doLog {
		logger.Println(cmd.Args)
	}
	cmd.Env = []string{
		"GOVC_INSECURE=" + helpers.Config.GovcInsecure,
		"GOVC_URL=" + helpers.Config.DynamicConfigs[baseEnv].GovcURL,
		"GOVC_USERNAME=" + helpers.Config.DynamicConfigs[baseEnv].GovcUsername,
		"GOVC_PASSWORD=" + helpers.Config.DynamicConfigs[baseEnv].GovcPassword,
		"GOVC_DATACENTER=" + helpers.Config.DynamicConfigs[baseEnv].GovcDatacenter,
		"GOVC_DATASTORE=" + helpers.Config.DynamicConfigs[baseEnv].GovcDatastore,
		"GOVC_RESOURCE_POOL=" + helpers.Config.DynamicConfigs[baseEnv].GovcResourcePool,
	}
	out, err := cmd.CombinedOutput()
	return out, err
}

func GetVMInfoDump(identifier string, targetName string, logger *log.Logger) ([]byte, error) {
	out, err := Execute(identifier, true, logger, "vm.info", "-json", targetName)
	return out, err
}
