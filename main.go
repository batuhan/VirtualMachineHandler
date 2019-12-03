package main

import (
	"github.com/nodgroup/VirtualMachineHandler/vmware"
	"net/http"
)

func main() {
	http.HandleFunc("/env", vmware.Env)
	http.HandleFunc("/create", vmware.Create)
	_ = http.ListenAndServe(":8080", nil)
}
