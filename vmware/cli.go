package vmware

import (
	"log"
	"os"
	"os/exec"
)

func execute(baseEnv string, doLog bool, logger *log.Logger, arg ...string) ([]byte, error) {
	cmd := exec.Command("govc", arg...)
	if doLog {
		logger.Println(cmd.Args)
	}
	cmd.Env = []string{
		"GOVC_INSECURE=" + os.Getenv("GOVC_INSECURE"),
		"GOVC_URL=" + os.Getenv(baseEnv+"_GOVC_URL"),
		"GOVC_USERNAME=" + os.Getenv(baseEnv+"_GOVC_USERNAME"),
		"GOVC_PASSWORD=" + os.Getenv(baseEnv+"_GOVC_PASSWORD"),
		"GOVC_DATACENTER=" + os.Getenv(baseEnv+"_GOVC_DATACENTER"),
		"GOVC_DATASTORE=" + os.Getenv(baseEnv+"_GOVC_DATASTORE"),
		"GOVC_RESOURCE_POOL=" + os.Getenv(baseEnv+"_GOVC_RESOURCE_POOL"),
	}
	out, err := cmd.CombinedOutput()
	return out, err
}
