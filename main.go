package main

import (
	"gitlab.com/nod/bigcore/VirtualMachineHandler/vmware"
	"net/http"
)

func main() {
	http.HandleFunc("/env", vmware.Env)
	http.HandleFunc("/create", func(w http.ResponseWriter, req *http.Request) {
		go vmware.Create(w, req)
	})
	_ = http.ListenAndServe(":8080", nil)
}
