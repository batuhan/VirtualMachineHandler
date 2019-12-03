package main

import (
	"github.com/nodgroup/VirtualMachineHandler/vmware"
	"net/http"
)

func main() {
	http.HandleFunc("/env", vmware.Env)
	_ = http.ListenAndServe(":8080", nil)
}
