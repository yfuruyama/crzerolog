crzerolog
===
[![godoc](https://godoc.org/github.com/yfuruyama/crzerolog?status.svg)](https://godoc.org/github.com/yfuruyama/crzerolog) [![CircleCI](https://circleci.com/gh/yfuruyama/crzerolog.svg?style=svg)](https://circleci.com/gh/yfuruyama/crzerolog)

A zerolog-based logging library for Cloud Run.

![screenshot](screenshot.png)

## Features

- Auto format with Stackdriver fields such as time, severity, trace, sourceLocation
- Group application logs with the request log
- Supports all of [rs/zerolog](https://github.com/rs/zerolog) APIs for structured logging

## Installation

```
go get -u github.com/yfuruyama/crzerolog
```

## Example

See [GoDoc](https://godoc.org/github.com/yfuruyama/crzerolog) for more details.

```
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
	middleware := crzerolog.InjectLogger(&rootLogger)

	http.Handle("/", middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := hlog.FromRequest(r)

		logger.Info().Msg("Hi")
		logger.Warn().Str("foo", "bar").Msg("This is")
		logger.Error().Int("num", 123).Msg("Structured Log")

		fmt.Fprintf(w, "Hello\n")
	})))

	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	log.Printf("Server listening on port %q", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
```

## Level mapping

This library automatically maps [zerolog level](https://godoc.org/github.com/rs/zerolog#Level) to [Stackdriver severity](https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogSeverity).

Mapping between them is as follows.

| zerolog level | Stackdriver severity |
| --- | --- |
| `NoLevel` | `DEFAULT` |
| `TraceLevel` | `DEFAULT` |
| `DebugLevel` | `DEBUG` |
| `InfoLevel` | `INFO` |
| `WarnLevel` | `WARNING` |
| `ErrorLevel` | `ERROR` |
| `FatalLevel` | `CRITICAL` |
| `PanicLevel` | `ALERT` |

## Supported Platform

- Cloud Run (fully managed)
- Google App Engine (2nd-Generation)

## License
[Apache 2.0](LICENSE).
