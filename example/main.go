package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/yfuruyama/crzerolog"
)

func main() {
	rootLogger := zerolog.New(os.Stdout)

	server := http.NewServeMux()
	server.Handle("/", crzerolog.HandleWithLogger(&rootLogger, http.HandlerFunc(hello)))

	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	log.Printf("Server listening on port %s", port)
	if err := http.ListenAndServe(":"+port, server); err != nil {
		log.Fatal(err)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	logger := hlog.FromRequest(r)

	logger.Info().Str("foo", "foo!").Msg("Hello structured log 1")
	logger.Warn().Str("bar", "bar!").Msg("Hello structured log 2")
	logger.Error().Str("baz", "baz!").Msg("Hello structured log 3")

	fmt.Fprintf(w, "Hello\n")
}
