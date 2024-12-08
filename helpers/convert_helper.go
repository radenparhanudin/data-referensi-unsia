package helpers

import (
	"regexp"
	"strings"
)

/* Convert Camel Case To Snake Case */
func ConvertCCToSC(str string) string {
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := re.ReplaceAllString(str, "${1}_${2}")
	return strings.ToLower(snake)
}
