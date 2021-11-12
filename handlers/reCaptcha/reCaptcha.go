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

	err = httputils.WriteJson(200, ReCaptchaResponse{Token: strings.TrimSuffix(stdOut.String(), "\n")}, w)
	if err != nil {
		log.Error(err)
		handlers.WriteHandlerError(err, http.StatusInternalServerError, httputils.UnexpectedError, w, r)
	}
}
