package handlers

import (
	"encoding/json"
	"net/http"
)

type HandlerError struct {
	Error string `json:"error"`
}

func WriteHandlerError(err string, code int, w http.ResponseWriter) {
	msgBytes, _ := json.Marshal(HandlerError{Error: err})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(msgBytes)
}

func WriteHandlerResp(resp interface{}, code int, w http.ResponseWriter) error {
	msgBytes, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(msgBytes)
	return nil
}

type ContactHandler interface {
	PostContact(w http.ResponseWriter, r *http.Request)
}

type DownloadHandler interface {
	GetResume(w http.ResponseWriter, r *http.Request)
}
type ReCaptchaHandler interface {
	GetToken(w http.ResponseWriter, r *http.Request)
}
type AnyToCsvHandler interface {
	PostCsv(w http.ResponseWriter, r *http.Request)
}
