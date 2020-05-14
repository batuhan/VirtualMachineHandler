package helpers

import (
	"encoding/json"
	"io"
	"log"
)

type Body struct {
	LocationId string
	TargetName string
}

type Delete struct {
	Body
}

type Update struct {
	Body
	Cpu      int
	Memory   int
	DiskSize string
}

type State struct {
	Body
	Action string
}

type Create struct {
	Body
	Template    string
	Cpu         int
	Memory      int
	DiskSize    string
	SshKeys     []string
	IpToAssign  string
	OnFirstBoot []string
}

func GetBody(reqBody io.ReadCloser) Body {
	body := Body{}
	err := json.NewDecoder(reqBody).Decode(&body)
	if err != nil {
		log.Println(err.Error())
	}
	return body
}
