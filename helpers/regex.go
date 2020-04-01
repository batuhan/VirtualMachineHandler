package helpers

import (
	"log"
	"regexp"
)

var CreateNameRegex *regexp.Regexp

func CompileRegexes() {
	createNameRegex, err := regexp.Compile(`\W|\s`)
	if err != nil {
		log.Fatal(err)
	}
	CreateNameRegex = createNameRegex
}

func ApplyCreateNameRegex(base string) string {
	return CreateNameRegex.ReplaceAllString(base, "")
}
