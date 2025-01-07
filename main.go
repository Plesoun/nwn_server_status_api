package main

import (
	"fmt"
	"net/http"
	"nwn_server_info/handlers"
	"nwn_server_info/logging"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func main() {
	logger, err := logging.NewLogger("stdout", "")
	if err != nil {
		fmt.Printf("failed to intialize logger: %s", err)
		os.Exit(1)
	}
	logger.WithField("UUID", uuid.New().String()).Info("Starting http server...")

	router := mux.NewRouter()

	router.HandleFunc("/players-online", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetPlayersOnline(logger, w, r)
	}).Methods("GET")

	if err := http.ListenAndServe(":8083", router); err != nil {
		logger.WithField("UUID", uuid.New().String()).Errorf("Failed to start server: %s", err)
		os.Exit(1)
	}
}
