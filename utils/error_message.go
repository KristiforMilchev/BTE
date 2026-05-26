package utils

import "fmt"

func ErrorMessage(err error) string {
	if err == nil {
		return "unknown error"
	}
	return fmt.Sprintf("%s", err)
}
