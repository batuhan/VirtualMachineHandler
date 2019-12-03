package vmware

import (
	"fmt"
	"net/http"
)

func Create(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprint(w, "test")
}
