package vmware

import (
	"fmt"
	"net/http"
)

func Env(w http.ResponseWriter, _ *http.Request) {
	out := execute("env")
	_, _ = fmt.Fprintf(w, out)
}
