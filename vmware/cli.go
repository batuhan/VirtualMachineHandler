package vmware

import (
	"log"
	"os/exec"
)

func execute(arg ...string) string {
	cmd := exec.Command("govc", arg...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with\n%s%s", out, err)
	}
	return string(out)
}
