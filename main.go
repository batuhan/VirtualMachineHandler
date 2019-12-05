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
		uuid := helpers.GenerateUUID()
		_, _ = fmt.Fprint(w, uuid)
		body := helpers.GetBody(req.Body)

		go vmware.Create(body, uuid)
	})
	http.HandleFunc("/delete", func(w http.ResponseWriter, req *http.Request) {
		uuid := helpers.GenerateUUID()
		_, _ = fmt.Fprint(w, uuid)
		body := helpers.GetBody(req.Body)

		go func() {
			_ = vmware.Delete(body, uuid)
		}()
	})
	http.HandleFunc("/recreate", func(w http.ResponseWriter, req *http.Request) {
		uuid := helpers.GenerateUUID()
		_, _ = fmt.Fprint(w, uuid)
		body := helpers.GetBody(req.Body)

		go vmware.Recreate(body, uuid)
	})
	http.HandleFunc("/update", func(w http.ResponseWriter, req *http.Request) {
		uuid := helpers.GenerateUUID()
		_, _ = fmt.Fprint(w, uuid)
		body := helpers.GetBody(req.Body)

		go vmware.Update(body, uuid)
	})
	http.HandleFunc("/dump", func(w http.ResponseWriter, req *http.Request) {
		dump, _ := httputil.DumpRequest(req, true)
		log.Println(string(dump))
	})
	_ = http.ListenAndServe(":8080", nil)
}
