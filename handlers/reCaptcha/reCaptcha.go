package recaptcha

import (
	"bytes"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/personal-website-backend/handlers"

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
		handlers.WriteHandlerError("no query parameters", http.StatusBadRequest, w)
		return
	}

	cmd := exec.Command("python3.6", "/app/applications/reCaptcha/app.py", url, clientKey, action)
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		log.Error(fmt.Sprint(err) + ": " + stdErr.String())
		handlers.WriteHandlerError("unable to get token", http.StatusInternalServerError, w)
		return
	}

	output := strings.TrimSuffix(stdOut.String(), "\n")
	switch output {
	case ERROR_RENDERING_RECAPTCHA:
		log.Error(ERROR_RENDERING_RECAPTCHA)
		handlers.WriteHandlerError(ERROR_RENDERING_RECAPTCHA, http.StatusBadRequest, w)
	default:
		err = handlers.WriteHandlerResp(ReCaptchaResponse{Token: output}, 200, w)
		if err != nil {
			log.Error(err)
			handlers.WriteHandlerError("unable to get token", http.StatusInternalServerError, w)
		}
	}
}
