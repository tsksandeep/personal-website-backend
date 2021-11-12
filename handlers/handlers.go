package handlers

import (
	"net/http"

	"github.com/website/httputils"
)

func WriteHandlerError(err error, httpErrorCode int, subErrorCode httputils.ErrorCode, w http.ResponseWriter, r *http.Request) {
	subError := httputils.NewSubError(subErrorCode, "message", err.Error())
	herr := httputils.NewHandlerError(httpErrorCode, subError)
	httputils.WriteHandlerError(herr, r, w)
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
