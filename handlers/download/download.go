package download

import (
	"net/http"
	"path/filepath"

	"github.com/website/handlers"

	log "github.com/sirupsen/logrus"
)

type downloadHanlder struct{}

// New creates a new instance of Contact Handler
func New() handlers.DownloadHandler {
	return &downloadHanlder{}
}

func (ch *downloadHanlder) GetResume(w http.ResponseWriter, r *http.Request) {
	absPath, err := filepath.Abs("./static/sandeep_resume.pdf")
	if err != nil {
		log.Error("error resolving relative path for sandeep resume")
		handlers.WriteHandlerError("internal server error", http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(absPath))
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, absPath)
}
