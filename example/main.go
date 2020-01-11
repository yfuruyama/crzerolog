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

	http.Handle("/", crzerolog.HandleWithLogger(&rootLogger, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := hlog.FromRequest(r)

		logger.Info().Str("foo", "foo!").Msg("Hello structured log 1")
		logger.Warn().Str("bar", "bar!").Msg("Hello structured log 2")
		logger.Error().Str("baz", "baz!").Msg("Hello structured log 3")

		fmt.Fprintf(w, "Hello\n")
	})))

	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	log.Printf("Server listening on port %q", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
