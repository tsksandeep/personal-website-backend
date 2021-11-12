package contact

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/website/email"
	"github.com/website/handlers"
	"github.com/website/httputils"

	log "github.com/sirupsen/logrus"
)

type contactHandler struct{}

// New creates a new instance of Contact Handler
func New() handlers.ContactHandler {
	return &contactHandler{}
}

func (ch *contactHandler) PostContact(w http.ResponseWriter, r *http.Request) {
	var contactInfo Contact

	err := json.NewDecoder(r.Body).Decode(&contactInfo)
	if err != nil {
		log.Error(err)
		handlers.WriteHandlerError(errors.New("no body parameter"), http.StatusBadRequest, httputils.BadRequest, w, r)
		return
	}

	err = email.SendEmail(&contactInfo)
	if err != nil {
		log.Error(err)
		handlers.WriteHandlerError(errors.New("email could not be sent"), http.StatusInternalServerError, httputils.UnexpectedError, w, r)
		return
	}

	err = httputils.WriteJson(200, nil, w)
	if err != nil {
		log.Error(err)
		handlers.WriteHandlerError(err, http.StatusInternalServerError, httputils.UnexpectedError, w, r)
	}
}
