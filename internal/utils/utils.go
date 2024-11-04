package utils

import (
	"regexp"
)

func ExtractSQLStateErrorCode(errMsg string) string {
	re := regexp.MustCompile(`SQLSTATE (\d{5})`)
	matches := re.FindStringSubmatch(errMsg)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

func SQLErrorMessage(statusCode string) string {
	switch statusCode {
	case "23505":
		return "email already in use"
	default:
		return "bad request"
	}
}
