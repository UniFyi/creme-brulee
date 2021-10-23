package messaging

import (
	"errors"
	"fmt"
)

var (
	ErrMissingField = errors.New("missing field")
)

type InvalidField struct {
	Name string
	Format string
}

func (i InvalidField) Error() string {
	return fmt.Sprintf("invalid field %v must be in %v format", i.Name, i.Format)
}

type ErrorResponse struct {
	Error ErrorResponsePayload `json:"error"`
}

type ErrorResponsePayload struct {
	Code    string `json:"code"`
	Details map[string]interface{} `json:"details"`
}

func CreateInvalidFieldError(fieldName string, message string) *ErrorResponse {
	return &ErrorResponse{Error: ErrorResponsePayload{
		Code:    "invalid-field",
		Details: map[string]interface{}{
			fieldName: message,
		},
	}}
}
func CreateUnsatisfiedRuleError(err error) *ErrorResponse {
	return &ErrorResponse{Error: ErrorResponsePayload{
		Code:    "rule-unsatisfied",
		Details: map[string]interface{}{
			"rule": err.Error(),
		},
	}}
}
func CreateNotFoundError(err error) *ErrorResponse {
	return &ErrorResponse{Error: ErrorResponsePayload{
		Code:    "not-found",
		Details: map[string]interface{}{
			"message": err.Error(),
		},
	}}
}
func CreateActionNotAllowedError(err error) *ErrorResponse {
	return &ErrorResponse{Error: ErrorResponsePayload{
		Code:    "action-not-permitted",
		Details: map[string]interface{}{
			"message": err.Error(),
		},
	}}
}
