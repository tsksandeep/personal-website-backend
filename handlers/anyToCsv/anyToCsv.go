package anytocsv

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/website/handlers"

	log "github.com/sirupsen/logrus"
)

var (
	PYTHON_FILE                   = "/app/applications/anyToCsv/app.py"
	CONVERSION_FILES_DIR          = "/app/applications/anyToCsv/csv-conversion-files/"
	UNABLE_TO_CONVERT_FILE        = "unable to convert file"
	ERROR_EXTENSION_NOT_SUPPORTED = "Exception: file extension not supported"
)

type anyToCsvHandler struct{}

// New creates a new instance of Contact Handler
func New() handlers.AnyToCsvHandler {
	return &anyToCsvHandler{}
}

func (ah *anyToCsvHandler) PostCsv(w http.ResponseWriter, r *http.Request) {
	var stdOut bytes.Buffer
	var stdErr bytes.Buffer

	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Error(err)
		return
	}

	log.Infof("file %s is being uploaded", handler.Filename)

	tempFile, err := ioutil.TempFile(CONVERSION_FILES_DIR, "*_"+handler.Filename)
	if err != nil {
		log.Error(err)
		return
	}

	defer file.Close()
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error(err)
		return
	}

	tempFile.Write(fileBytes)

	log.Infof("%s uploaded", handler.Filename)

	cmd := exec.Command("python3.6", PYTHON_FILE, tempFile.Name())
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr
	err = cmd.Run()
	if err != nil {
		errOutput := strings.TrimSuffix(stdErr.String(), "\n")
		switch errOutput {
		case ERROR_EXTENSION_NOT_SUPPORTED:
			log.Error(errOutput)
			handlers.WriteHandlerError("file extension not supported", http.StatusBadRequest, w)
		default:
			log.Error(errOutput)
			handlers.WriteHandlerError(UNABLE_TO_CONVERT_FILE, http.StatusInternalServerError, w)
		}
		return
	}

	output := strings.TrimSuffix(stdOut.String(), "\n")
	log.Infof("successfully converted %s", output)
	defer os.Remove(output)

	w.WriteHeader(200)
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(output))
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, output)
}
