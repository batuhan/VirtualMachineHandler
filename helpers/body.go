package helpers

import (
	"encoding/json"
	"io"
	"log"
)

type Body struct {
	Identifier  string
	Template    string
	TargetName  string
	Cpu         int
	Memory      int
	DiskSize    string
	SshKey      string
	IpToAssign  string
	Action      string
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
