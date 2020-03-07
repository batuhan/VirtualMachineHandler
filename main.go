package main

import (
	"encoding/json"
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
		body := helpers.Create{}
		err := json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			log.Println(err.Error())
		}
		go vmware.Create(body, uuid)
	})
	http.HandleFunc("/delete", func(w http.ResponseWriter, req *http.Request) {
		uuid := helpers.GenerateUUID()
		_, _ = fmt.Fprint(w, uuid)
		body := helpers.Delete{}
		err := json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			log.Println(err.Error())
		}
		go func() {
			_ = vmware.Delete(body.Identifier, body.TargetName, uuid)
		}()
	})
	http.HandleFunc("/recreate", func(w http.ResponseWriter, req *http.Request) {
		uuid := helpers.GenerateUUID()
		_, _ = fmt.Fprint(w, uuid)
		body := helpers.Create{}
		err := json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			log.Println(err.Error())
		}
		go vmware.Recreate(body, uuid)
	})
	http.HandleFunc("/update", func(w http.ResponseWriter, req *http.Request) {
		uuid := helpers.GenerateUUID()
		_, _ = fmt.Fprint(w, uuid)
		body := helpers.Update{}
		err := json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			log.Println(err.Error())
		}
		go vmware.Update(body, uuid)
	})
	http.HandleFunc("/state", func(w http.ResponseWriter, req *http.Request) {
		uuid := helpers.GenerateUUID()
		_, _ = fmt.Fprint(w, uuid)
		body := helpers.State{}
		err := json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			log.Println(err.Error())
		}
		go vmware.State(body, uuid)
	})

	http.HandleFunc("/dump", func(w http.ResponseWriter, req *http.Request) {
		dump, _ := httputil.DumpRequest(req, true)
		log.Println(string(dump))
	})
	_ = http.ListenAndServe(":8080", nil)
}
