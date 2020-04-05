package helpers

import (
	"log"
	"regexp"
)

var TargetNameRegex *regexp.Regexp

func CompileRegexes() {
	targetNameRegex, err := regexp.Compile(`\W|\s`)
	if err != nil {
		log.Fatal(err)
	}
	TargetNameRegex = targetNameRegex
}

func ApplyTargetNameRegex(base string) string {
	return TargetNameRegex.ReplaceAllString(base, "")
}
