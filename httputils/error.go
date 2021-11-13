package httputils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var ErrorCodeStrings = map[ErrorCode]string{
	FormatError:         "format_error",
	NotFound:            "not_found",
	BadRequest:          "bad_request",
	InvalidScope:        "invalid_scope",
	UnexpectedError:     "unexpected_server_error",
	NotImplemented:      "not_implemented",
	InvalidOperation:    "invalid_operation",
	InvalidParameter:    "invalid_parameter",
	Deprecated:          "deprecated",
	Forbidden:           "forbidden",
	PreconditionFailed:  "precondition_failed",
	UnprocessableEntity: "unprocessable_entity",
}

type ErrorCode int

const (
	_ ErrorCode = iota
	Custom
	NotFound
	FormatError
	BadRequest
	InvalidScope
	UnexpectedError
	NotImplemented
	InvalidOperation
	InvalidParameter
	Deprecated
	Forbidden
	PreconditionFailed
	UnprocessableEntity
)

type ErrorDetails map[string]interface{}

type SubError struct {
	Code    ErrorCode
	Details ErrorDetails
}

func (subError *SubError) MarshalJSON() ([]byte, error) {
	if subError.Code != Custom {
		subError.Details["code"] = ErrorCodeStrings[subError.Code]
	}
	return json.Marshal(subError.Details)
}

func (subError *SubError) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Code: %s", ErrorCodeStrings[subError.Code])
	for key, value := range subError.Details {
		fmt.Fprintf(&buf, "\n%s: %s", key, value)
	}
	return buf.String()
}

func NewSubError(code ErrorCode, key string, value interface{}) *SubError {
	return &SubError{
		Code:    code,
		Details: ErrorDetails{key: value},
	}
}

type HandlerError struct {
	HttpStatusCode int
	SubErrors      []*SubError
}

func (handlerErr *HandlerError) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "HttpStatusCode: %d", handlerErr.HttpStatusCode)
	for _, subError := range handlerErr.SubErrors {
		fmt.Fprintf(&buf, "\n%s", subError)
	}
	return buf.String()
}

func NewHandlerError(statusCode int, subErrors ...*SubError) *HandlerError {
	return &HandlerError{
		HttpStatusCode: statusCode,
		SubErrors:      subErrors,
	}
}

func NewNotFoundError(message string) *HandlerError {
	subError := NewSubError(NotFound, "message", message)
	return NewHandlerError(http.StatusNotFound, subError)
}
func NewFormatError(message string) *HandlerError {
	subError := NewSubError(FormatError, "message", message)
	return NewHandlerError(http.StatusBadRequest, subError)
}
func NewInvalidOperation(message string) *HandlerError {
	subError := NewSubError(InvalidOperation, "message", message)
	return NewHandlerError(http.StatusConflict, subError)
}
func NewInvalidParameterError(message string) *HandlerError {
	subError := NewSubError(InvalidParameter, "message", message)
	return NewHandlerError(http.StatusBadRequest, subError)
}
func NewUnexpectedError(err error) *HandlerError {
	subError := NewSubError(UnexpectedError, "error", err.Error())
	return NewHandlerError(http.StatusInternalServerError, subError)
}
func NewDeprecatedError(message string) *HandlerError {
	subError := NewSubError(Deprecated, "message", message)
	return NewHandlerError(http.StatusGone, subError)
}

func NewCustomError(httpStatus int, code, message string) *HandlerError {
	subError := NewSubError(Custom, "code", code)
	subError.Details["message"] = message
	return NewHandlerError(httpStatus, subError)
}
func NewBadRequestError(errors ErrorDetails) *HandlerError {
	subError := SubError{
		Code:    BadRequest,
		Details: errors,
	}
	return NewHandlerError(http.StatusBadRequest, &subError)
}
func NewForbiddenError(message string) *HandlerError {
	subError := NewSubError(Forbidden, "message", message)
	return NewHandlerError(http.StatusForbidden, subError)
}

type HttpError struct {
	Status    int         `json:"httpStatus"`
	Code      string      `json:"httpCode"`
	Errors    []*SubError `json:"errors"`
}
