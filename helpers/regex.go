package helpers

import (
	"log"
	"regexp"
)

var TargetNameRegex *regexp.Regexp

func CompileRegexes() {
	createNameRegex, err := regexp.Compile(`\W|\s`)
	if err != nil {
		log.Fatal(err)
	}
	TargetNameRegex = createNameRegex
}

func ApplyTargetNameRegex(base string) string {
	return TargetNameRegex.ReplaceAllString(base, "")
}
