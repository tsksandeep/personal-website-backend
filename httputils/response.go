package httputils

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

var ContextRequestIDKey interface{} = "requestId"

var httpStatusCodes = map[int]string{
	http.StatusInternalServerError: "internal_server_error",
	http.StatusConflict:            "conflict",
	http.StatusNotFound:            "not_found",
	http.StatusBadRequest:          "bad_request",
	http.StatusUnauthorized:        "unauthorized",
	http.StatusForbidden:           "forbidden",
	http.StatusUnprocessableEntity: "unprocessable_entity",
	http.StatusPreconditionFailed:  "precondition_failed",
}

var IsFatalError func(err error) bool

type HandlerFunc func(http.ResponseWriter, *http.Request, httprouter.Params) *HandlerError

func ToHandle(h HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		handlerErr := h(w, r, p)
		if handlerErr != nil {
			WriteHandlerError(handlerErr, r, w)
		}
	}
}

func AbbreAuthToken(authToken string) string {
	charsToReveal := 4
	if len(authToken) < charsToReveal {
		charsToReveal = len(authToken)
	}
	return authToken[:charsToReveal] + "..."
}

func GetAbbreAuthToken(r *http.Request) string {
	authToken := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)
	return AbbreAuthToken(authToken)
}

func WriteHandlerError(handlerErr *HandlerError, r *http.Request, w http.ResponseWriter) {
	httpError := &HttpError{
		Status: handlerErr.HttpStatusCode,
		Errors: handlerErr.SubErrors,
		Code:   httpStatusCodes[handlerErr.HttpStatusCode],
	}

	logFields := map[string]interface{}{
		"error":      handlerErr,
		"requestURI": r.RequestURI,
		"method":     r.Method,
		"authToken":  GetAbbreAuthToken(r),
	}
	log.WithFields(logFields).Error("request failed")

	if handlerErr.HttpStatusCode == http.StatusInternalServerError ||
		handlerErr.HttpStatusCode == http.StatusForbidden ||
		handlerErr.HttpStatusCode == http.StatusUnauthorized {
		httpError.Errors = make([]*SubError, 0)
	}

	err := WriteJson(handlerErr.HttpStatusCode, httpError, w)
	if err != nil {
		log.Errorf("serializing http error failed: %v", err.Error())
	}

	if handlerErr.HttpStatusCode == http.StatusInternalServerError && isFatalError(handlerErr) {
		log.Fatal("fatal error in the application")
	}
}

func isFatalError(handlerErr *HandlerError) bool {
	if IsFatalError == nil {
		return false
	}
	for _, subError := range handlerErr.SubErrors {
		if err, ok := subError.Details["error"].(error); ok {
			return IsFatalError(err)
		}
	}
	return false
}

func WriteJson(status int, response interface{}, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	if response == nil {
		return nil
	}

	return json.NewEncoder(w).Encode(response)
}
