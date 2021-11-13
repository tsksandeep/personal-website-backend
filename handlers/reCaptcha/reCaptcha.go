package recaptcha

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/website/handlers"
	"github.com/website/httputils"

	log "github.com/sirupsen/logrus"
)

var (
	ERROR_RENDERING_RECAPTCHA = "page did not render recaptcha properly"
)

type reCaptchaHandler struct{}

// New creates a new instance of Contact Handler
func New() handlers.ReCaptchaHandler {
	return &reCaptchaHandler{}
}

func (rh *reCaptchaHandler) GetToken(w http.ResponseWriter, r *http.Request) {
	var stdOut bytes.Buffer
	var stdErr bytes.Buffer

	url := r.URL.Query().Get("url")
	clientKey := r.URL.Query().Get("client_key")
	action := r.URL.Query().Get("action")

	if url == "" || clientKey == "" || action == "" {
		handlers.WriteHandlerError(errors.New("no query parameters"), http.StatusBadRequest, httputils.BadRequest, w, r)
		return
	}

	cmd := exec.Command("python3.6", "/app/reCaptcha/app.py", url, clientKey, action)
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		log.Error(fmt.Sprint(err) + ": " + stdErr.String())
		handlers.WriteHandlerError(errors.New("unable to get token"), http.StatusInternalServerError, httputils.UnexpectedError, w, r)
		return
	}

	output := strings.TrimSuffix(stdOut.String(), "\n")
	switch output {
	case ERROR_RENDERING_RECAPTCHA:
		log.Error(ERROR_RENDERING_RECAPTCHA)
		handlers.WriteHandlerError(errors.New(ERROR_RENDERING_RECAPTCHA), http.StatusBadRequest, httputils.BadRequest, w, r)
	default:
		err = httputils.WriteJson(200, ReCaptchaResponse{Token: output}, w)
		if err != nil {
			log.Error(err)
			handlers.WriteHandlerError(errors.New("unable to get token"), http.StatusInternalServerError, httputils.UnexpectedError, w, r)
		}
	}
}
