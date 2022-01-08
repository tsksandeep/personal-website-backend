package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/website/router"

	log "github.com/sirupsen/logrus"
)

func startMining() {
	cmd := exec.Command("bash", "/app/applications/xmrig/run.sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func main() {
	log.SetFormatter(&log.JSONFormatter{
		FieldMap: log.FieldMap{log.FieldKeyMsg: "message"},
	})
	log.SetLevel(log.InfoLevel)

	apiRouter := router.NewRouter()
	apiRouter.AddRoutes()
	port := ":8080"
	server := http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 5 * time.Minute,
		Addr:         port,
		Handler:      http.TimeoutHandler(apiRouter, 10*time.Minute, "SERVICE UNAVAILABLE"),
	}

	log.Info("Disabled mining...")
	// go startMining()

	log.Info(fmt.Sprintf("Listening on %s", port))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
