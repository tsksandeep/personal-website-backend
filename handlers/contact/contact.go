package contact

import (
	"encoding/json"
	"net/http"

	"github.com/personal-website-backend/email"
	"github.com/personal-website-backend/handlers"

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
		handlers.WriteHandlerError("no body parameter", http.StatusBadRequest, w)
		return
	}

	err = email.SendEmail(&contactInfo)
	if err != nil {
		log.Error(err)
		handlers.WriteHandlerError("email could not be sent", http.StatusInternalServerError, w)
		return
	}

	err = handlers.WriteHandlerResp(nil, http.StatusOK, w)
	if err != nil {
		log.Error(err)
		handlers.WriteHandlerError("internal server error", http.StatusInternalServerError, w)
	}
}
