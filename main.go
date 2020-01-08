package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

func main() {
	port := "8080"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}

	rootLogger := zerolog.New(os.Stdout)

	server := http.NewServeMux()
	server.Handle("/", HandleWithLogger(&rootLogger, http.HandlerFunc(hello)))
	log.Printf("Server listening on port %s", port)
	err := http.ListenAndServe(":"+port, server)
	log.Fatal(err)
}

func hello(w http.ResponseWriter, r *http.Request) {
	logger := hlog.FromRequest(r)
	logger.Info().Str("additionalField", "test").Msg("Hello structured log 1")
	logger.Warn().Str("additionalField", "test").Msg("Hello structured log 2")

	status := r.URL.Query().Get("status")
	statusCode, err := strconv.Atoi(status)
	if status == "" || err != nil {
		statusCode = 200
	}

	w.WriteHeader(statusCode)
	fmt.Fprintf(w, "Hello! %s", http.StatusText(statusCode))
}
