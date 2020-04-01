package main

import (
	"encoding/json"
	"fmt"
	"gitlab.com/nod/bigcore/VirtualMachineHandler/helpers"
	"gitlab.com/nod/bigcore/VirtualMachineHandler/vmware/actions"
	"log"
	"net/http"
	"net/http/httputil"
)

func main() {
	helpers.InitConfig()
	helpers.CompileRegexes()

	http.HandleFunc("/env", actions.Env)
	http.HandleFunc("/create", func(w http.ResponseWriter, req *http.Request) {
		uuid := helpers.GenerateUUID()
		_, _ = fmt.Fprint(w, uuid)
		body := helpers.Create{}
		err := json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			log.Println(err.Error())
		}
		go actions.Create(body, uuid)
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
			_ = actions.Delete(body.Identifier, body.TargetName, uuid)
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
		go actions.Recreate(body, uuid)
	})
	http.HandleFunc("/update", func(w http.ResponseWriter, req *http.Request) {
		uuid := helpers.GenerateUUID()
		_, _ = fmt.Fprint(w, uuid)
		body := helpers.Update{}
		err := json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			log.Println(err.Error())
		}
		go actions.Update(body, uuid)
	})
	http.HandleFunc("/changeState", func(w http.ResponseWriter, req *http.Request) {
		uuid := helpers.GenerateUUID()
		_, _ = fmt.Fprint(w, uuid)
		body := helpers.State{}
		err := json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			log.Println(err.Error())
		}
		go actions.ChangeState(body, uuid)
	})

	http.HandleFunc("/dump", func(w http.ResponseWriter, req *http.Request) {
		dump, _ := httputil.DumpRequest(req, true)
		log.Println(string(dump))
	})

	log.Printf("server is listening at port %s", helpers.Config.HttpPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", helpers.Config.HttpPort), nil))
}
