package main

import (
	"github.com/nodgroup/VirtualMachineHandler/vmware"
	"net/http"
)

func main() {
	http.HandleFunc("/create_server", vmware.Create)
	_ = http.ListenAndServe(":8080", nil)
}
