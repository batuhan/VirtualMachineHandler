package helpers

import (
	"encoding/json"
	"io"
)

type Body struct {
	Identifier string
	TargetName string
}

func GetBody(reqBody io.ReadCloser) Body {
	body := Body{}
	_ = json.NewDecoder(reqBody).Decode(&body)
	return body
}
