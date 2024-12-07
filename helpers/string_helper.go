package helpers

import (
	"regexp"
	"strings"
)

func CamelCaseToSnakeCase(str string) string {
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := re.ReplaceAllString(str, "${1}_${2}")
	return strings.ToLower(snake)
}

func GenerateValidationErrorMessage(fieldName string, tag string) string {
	errorMessages := map[string]string{
		"required":  fieldName + " is required.",
		"numeric":   fieldName + " must be a number.",
		"omitempty": fieldName + " is optional.",
	}

	if message, exists := errorMessages[tag]; exists {
		return message
	}
	return fieldName + " is invalid."
}
