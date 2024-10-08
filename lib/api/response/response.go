package response

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

func Success() Response {
	return Response{Status: StatusSuccess}
}
func Fail(message string) Response {
	return Response{
		Status: StatusError,
		Error:  message,
	}
}
func ValidationError(errs validator.ValidationErrors) Response {
	var errorMessages []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errorMessages = append(errorMessages, fmt.Sprintf("%s is required", err.Field()))
		case "url":
			errorMessages = append(errorMessages, fmt.Sprintf("%s is not a valid URL", err.Field()))
		default:
			errorMessages = append(errorMessages, fmt.Sprintf("%s is invalid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errorMessages, ", "),
	}
}
