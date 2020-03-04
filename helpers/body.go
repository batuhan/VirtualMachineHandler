package helpers

import (
	"encoding/json"
	"io"
	"log"
)

type Body struct {
	Identifier string
}

type Delete struct {
	Body
	TargetName string
}

type Update struct {
	Body
	TargetName string
	Cpu        int
	Memory     int
	DiskSize   string
}

type State struct {
	Body
	TargetName string
	Action     string
}

type Create struct {
	Body
	Template    string
	TargetName  string
	Cpu         int
	Memory      int
	DiskSize    string
	SshKey      string
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
