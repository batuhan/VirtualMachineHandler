package helpers

import (
	"encoding/json"
	"io"
)

type Body struct {
	Identifier string
	TargetName string
	Cpu        int
	Memory     int
}

func GetBody(reqBody io.ReadCloser) Body {
	body := Body{}
	_ = json.NewDecoder(reqBody).Decode(&body)
	return body
}
