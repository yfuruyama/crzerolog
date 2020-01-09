package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/yfuruyama/cloudrunlog"
)

func main() {
	port := "8080"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}

	rootLogger := zerolog.New(os.Stdout)

	server := http.NewServeMux()
	server.Handle("/", cloudrunlog.HandleWithLogger(&rootLogger, http.HandlerFunc(hello)))
	log.Printf("Server listening on port %s", port)
	err := http.ListenAndServe(":"+port, server)
	log.Fatal(err)
}

func hello(w http.ResponseWriter, r *http.Request) {
	logger := hlog.FromRequest(r)
	logger.Info().Str("foo", "foo!").Msg("Hello structured log 1")
	logger.Warn().Str("bar", "bar!").Msg("Hello structured log 2")
	logger.Error().Str("baz", "baz!").Msg("Hello structured log 3")

	status := r.URL.Query().Get("status")
	statusCode, err := strconv.Atoi(status)
	if status == "" || err != nil {
		statusCode = 200
	}

	w.WriteHeader(statusCode)
	fmt.Fprintf(w, "Hello! %s\n", http.StatusText(statusCode))
}
