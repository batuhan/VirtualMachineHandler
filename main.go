package main

import (
	"fmt"
	"gitlab.com/nod/bigcore/VirtualMachineHandler/helpers"
	"gitlab.com/nod/bigcore/VirtualMachineHandler/vmware"
	"log"
	"net/http"
	"net/http/httputil"
)

func main() {
	http.HandleFunc("/env", vmware.Env)
	http.HandleFunc("/create", func(w http.ResponseWriter, req *http.Request) {
		newUUID := helpers.GenerateUUID()
		_, _ = fmt.Fprint(w, newUUID)
		body := helpers.GetBody(req.Body)

		go vmware.Create(body, newUUID)
	})
	http.HandleFunc("/dump", func(w http.ResponseWriter, req *http.Request) {
		dump, _ := httputil.DumpRequest(req, true)
		log.Println(string(dump))
	})
	_ = http.ListenAndServe(":8080", nil)
}
