package helpers

import (
	"log"
	"regexp"
)

var CreateNameRegex *regexp.Regexp

func CompileRegexes() {
	var createNameRegex, err = regexp.Compile(`\W|\s`)
	if err != nil {
		log.Fatal(err)
	}
	CreateNameRegex = createNameRegex
}
